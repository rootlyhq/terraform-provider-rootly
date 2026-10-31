package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func TestAccResourceApiKey(t *testing.T) {
	if os.Getenv("ROOTLY_TEST_API_KEY") == "" {
		t.Skip("Skipping: CI uses a service account token which cannot create personal API keys")
	}
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyConfig(rName, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
				),
			},
		},
	})
}

func TestAccResourceApiKeyTeamScoped(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)
	var firstID string

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyTeamScopedConfig(rName, expiresAt, "first"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_api_key.test", "kind", "team"),
					resource.TestCheckResourceAttrPair("rootly_api_key.test", "group_id", "rootly_team.first", "id"),
					testAccCaptureResourceID("rootly_api_key.test", &firstID),
				),
			},
			{
				ResourceName:            "rootly_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			{
				Config: testAccResourceApiKeyTeamScopedConfig(rName, expiresAt, "second"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("rootly_api_key.test", "group_id", "rootly_team.second", "id"),
					testAccCheckResourceIDChanged("rootly_api_key.test", &firstID),
				),
			},
		},
	})
}

func TestAccResourceApiKeyRotation(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyRotationConfig(rName, expiresAt, "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_api_key.test", "token"),
					resource.TestCheckResourceAttr("rootly_api_key.test", "rotation_trigger", "1"),
				),
			},
			{
				ResourceName:            "rootly_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "rotation_trigger"},
			},
			// The rotate endpoint only lets a key rotate itself, and the test provider
			// authenticates with a different key. TestApiKeyRotateOnWire covers the
			// success path against a mock server.
			{
				Config:      testAccResourceApiKeyRotationConfig(rName, expiresAt, "2"),
				ExpectError: regexp.MustCompile(`Error rotating api_key: .*\(403\)`),
			},
		},
	})
}

func testAccCaptureResourceID(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckResourceIDChanged(name string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		if rs.Primary.ID == *previous {
			return fmt.Errorf("expected %s to be replaced, but ID %s is unchanged", name, rs.Primary.ID)
		}
		return nil
	}
}

func testAccResourceApiKeyConfig(name, expiresAt string) string {
	return fmt.Sprintf(`
resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "%s"
}
`, name, expiresAt)
}

func testAccResourceApiKeyTeamScopedConfig(name, expiresAt, team string) string {
	return fmt.Sprintf(`
resource "rootly_team" "first" {
	name = "%[1]s-first"
}

resource "rootly_team" "second" {
	name = "%[1]s-second"
}

resource "rootly_api_key" "test" {
	name       = "%[1]s"
	expires_at = "%[2]s"
	kind       = "team"
	group_id   = rootly_team.%[3]s.id
}
`, name, expiresAt, team)
}

func testAccResourceApiKeyRotationConfig(name, expiresAt, version string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%[1]s"
}

resource "rootly_api_key" "test" {
	name       = "%[1]s"
	expires_at = "%[2]s"
	kind       = "team"
	group_id   = rootly_team.test.id

	rotation_trigger = "%[3]s"
}
`, name, expiresAt, version)
}

func TestApiKeyRotateOnWire(t *testing.T) {
	var requests []string
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/vnd.api+json")
		if r.Method == http.MethodPost && r.URL.Path == "/v1/api_keys/1/rotate" {
			var body struct {
				Data struct {
					Attributes map[string]interface{} `json:"attributes"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decoding rotate request body: %v", err)
			}
			got = body.Data.Attributes
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id": "1", "type": "api_keys", "attributes": map[string]interface{}{
					"name":                 "deploy",
					"kind":                 "personal",
					"expires_at":           "2027-01-01T00:00:00Z",
					"grace_period_ends_at": "2026-10-26T00:10:00Z",
					"token":                "new-token",
				},
			},
		})
	}))
	t.Cleanup(server.Close)

	api, err := client.NewClient(server.URL, "test-token", "api-key-test")
	if err != nil {
		t.Fatal(err)
	}

	r := resourceApiKey()
	state := &terraform.InstanceState{ID: "1", Attributes: map[string]string{
		"id":               "1",
		"name":             "deploy",
		"kind":             "personal",
		"expires_at":       "2027-01-01T00:00:00Z",
		"token":            "old-token",
		"rotation_trigger": "1",
	}}
	diff := &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"rotation_trigger":              {Old: "1", New: "2"},
			"rotation_grace_period_minutes": {Old: "", New: "0"},
			"token":                         {Old: "old-token", NewComputed: true},
		},
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"rotation_grace_period_minutes": cty.NumberIntVal(0),
		}),
	}
	data, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatal(err)
	}

	if diags := resourceApiKeyUpdate(context.Background(), data, api); diags.HasError() {
		t.Fatal(diags)
	}

	if len(requests) != 1 || requests[0] != "POST /v1/api_keys/1/rotate" {
		t.Fatalf("expected only the rotate request, got %v", requests)
	}
	if got["expires_at"] != "2027-01-01T00:00:00Z" {
		t.Errorf("expected expires_at to be sent, got %#v", got["expires_at"])
	}
	if got["grace_period_minutes"] != float64(0) {
		t.Errorf("expected grace_period_minutes 0, got %#v", got["grace_period_minutes"])
	}
	if token := data.Get("token").(string); token != "new-token" {
		t.Errorf("expected token from rotate response, got %q", token)
	}
	if ends := data.Get("grace_period_ends_at").(string); ends != "2026-10-26T00:10:00Z" {
		t.Errorf("expected grace_period_ends_at from rotate response, got %q", ends)
	}
}
