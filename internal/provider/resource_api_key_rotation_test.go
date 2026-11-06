package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/jsonapi"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/samber/lo"
)

func TestAccResourceApiKeyRotation_Basic(t *testing.T) {
	addr := "rootly_api_key_rotation.test"

	// The resource can only rotate the API key the provider is authenticated with. Rotate a throwaway key instead of
	// the key the tests run with.
	apiKey := testAccCreateApiKey(t, acctest.RandomWithPrefix("tf-api-key-rotation"))
	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour).Truncate(time.Second).Format(time.RFC3339)

	tokenDiffers := statecheck.CompareValue(compare.ValuesDiffer())

	// The provider reads its token from a file, so that a step can switch it to the token issued by the previous
	// rotation, like a real rotation workflow would.
	tokenFile := filepath.Join(t.TempDir(), "token")
	var lastToken string

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccWriteTokenFile(t, tokenFile, &apiKey.Token),
				Config:    testAccResourceApiKeyRotationConfig(tokenFile, apiKey.ID, "1", expiresAt),
				Check:     testAccExtractToken(addr, &lastToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("api_key_id"), knownvalue.StringExact(apiKey.ID)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("grace_period_minutes"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("token_expires_at"), knownvalue.NotNull()),
					tokenDiffers.AddStateValue(addr, tfjsonpath.New("token")),
				},
			},
			{
				// Switch the provider to the rotated token, then change the triggers to rotate the API key again. The
				// grace period only covers the token that was just replaced, so the original token no longer works.
				PreConfig: testAccWriteTokenFile(t, tokenFile, &lastToken),
				Config:    testAccResourceApiKeyRotationConfig(tokenFile, apiKey.ID, "2", expiresAt),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(addr, plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					tokenDiffers.AddStateValue(addr, tfjsonpath.New("token")),
				},
			},
		},
	})
}

func testAccResourceApiKeyRotationConfig(tokenFile, apiKeyID, trigger, expiresAt string) string {
	return fmt.Sprintf(`
provider "rootly" {
	api_token = trimspace(file(%[1]q))
}

resource "rootly_api_key_rotation" "test" {
	api_key_id = %[2]q
	expires_at = %[4]q

	triggers = {
		rotation = %[3]q
	}
}
`, tokenFile, apiKeyID, trigger, expiresAt)
}

func testAccWriteTokenFile(t *testing.T, path string, token *string) func() {
	return func() {
		if err := os.WriteFile(path, []byte(*token), 0o600); err != nil {
			t.Fatalf("unable to write token file: %s", err)
		}
	}
}

func testAccExtractToken(addr string, token *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("not found in state: %s", addr)
		}
		*token = rs.Primary.Attributes["token"]
		return nil
	}
}

// testAccCreateApiKey creates a personal API key for the test user, and revokes it when the test finishes.
func testAccCreateApiKey(t *testing.T, name string) *apiclient.ApiKey {
	t.Helper()

	ctx := context.Background()
	apiKey, err := acctest.SharedClient.CreateApiKey(ctx, apiclient.ApiKey{
		Name:      name,
		Kind:      "personal",
		ExpiresAt: jsonapi.NewNullableAttrWithValue(time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)),
	})
	if err != nil {
		t.Fatalf("unable to create API key: %s", err)
	}

	t.Cleanup(func() {
		if err := acctest.SharedClient.DeleteApiKey(ctx, apiKey.ID); err != nil && !errors.Is(err, client.NotFoundError{}) {
			t.Errorf("unable to delete API key %s: %s", apiKey.ID, err)
		}
	})

	return apiKey
}

func TestResourceApiKeyRotation_ValidateConfig(t *testing.T) {
	ctx := context.Background()
	r := NewApiKeyRotationResource().(*ApiKeyRotationResource)

	schemaResp := &fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, schemaResp)
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	newConfig := func(gracePeriodMinutes *int64, expiresAt *string) tfsdk.Config {
		attrs := map[string]tftypes.Value{}
		for name, typ := range objectType.AttributeTypes {
			attrs[name] = tftypes.NewValue(typ, nil)
		}
		attrs["api_key_id"] = tftypes.NewValue(tftypes.String, "api-key-id")
		if gracePeriodMinutes != nil {
			attrs["grace_period_minutes"] = tftypes.NewValue(tftypes.Number, *gracePeriodMinutes)
		}
		if expiresAt != nil {
			attrs["expires_at"] = tftypes.NewValue(tftypes.String, *expiresAt)
		}
		return tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    tftypes.NewValue(objectType, attrs),
		}
	}

	testCases := map[string]struct {
		gracePeriodMinutes *int64
		expiresAt          *string
		expectWarnings     int
		expectErrors       int
	}{
		"defaults":           {},
		"grace period":       {gracePeriodMinutes: lo.ToPtr[int64](60)},
		"zero grace period":  {gracePeriodMinutes: lo.ToPtr[int64](0), expectWarnings: 1},
		"valid expires_at":   {expiresAt: lo.ToPtr("2027-01-01T00:00:00Z")},
		"invalid expires_at": {expiresAt: lo.ToPtr("2027-01-01"), expectErrors: 1},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resp := &fwresource.ValidateConfigResponse{}
			r.ValidateConfig(ctx, fwresource.ValidateConfigRequest{Config: newConfig(tc.gracePeriodMinutes, tc.expiresAt)}, resp)

			if got := resp.Diagnostics.WarningsCount(); got != tc.expectWarnings {
				t.Errorf("expected %d warnings, got %d: %v", tc.expectWarnings, got, resp.Diagnostics)
			}
			if got := resp.Diagnostics.ErrorsCount(); got != tc.expectErrors {
				t.Errorf("expected %d errors, got %d: %v", tc.expectErrors, got, resp.Diagnostics)
			}
		})
	}
}
