package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

func init() {
	resource.AddTestSweepers("rootly_live_call_router", &resource.Sweeper{
		Name: "rootly_live_call_router",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListLiveCallRoutersParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListLiveCallRoutersWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting live call routers, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting live call routers, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting live call routers, got empty response")
				}

				for _, liveCallRouter := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(liveCallRouter.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteLiveCallRouterWithResponse(ctx, liveCallRouter.Id)
						if err != nil {
							return fmt.Errorf("Error deleting live call router: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting live call router, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted live call router %s", liveCallRouter.Attributes.Name)
					}
				}

				if httpResp.ApplicationvndApiJSON200.Links.Next == nil {
					break
				}

				params.PageNumber = ptr.Ptr(ptr.Value(params.PageNumber) + 1)
			}

			return nil
		},
	})
}

func generatePhoneNumberForLiveCallRouter(ctx context.Context) (string, error) {
	httpResp, err := acctest.SharedClient.GeneratePhoneNumberLiveCallRouterWithResponse(ctx, &rootly.GeneratePhoneNumberLiveCallRouterParams{
		CountryCode: "US",
		PhoneType:   "local",
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate phone number: %v", err)
	}

	var generatePhoneNumberResp struct {
		PhoneNumber string `json:"phone_number"`
	}
	err = json.Unmarshal(httpResp.Body, &generatePhoneNumberResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal phone number: %v", err)
	}

	if generatePhoneNumberResp.PhoneNumber == "" {
		return "", fmt.Errorf("phone number is empty")
	}

	return generatePhoneNumberResp.PhoneNumber, nil
}

func TestAccResourceLiveCallRouter(t *testing.T) {
	resName := "rootly_live_call_router.test"
	name := acctest.RandomWithPrefix("tf-live-call-router")

	phoneNumber, err := generatePhoneNumberForLiveCallRouter(t.Context())
	if err != nil {
		t.Fatalf("Failed to generate phone number: %v", err)
	}

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("country_code"), knownvalue.StringExact("US")),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("phone_type"), knownvalue.StringExact("local")),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("phone_number"), knownvalue.StringExact(phoneNumber)),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("voicemail_greeting"), knownvalue.StringExact("Thank you for calling Rootly.")),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("caller_greeting"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("waiting_music_url"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("paging_targets"), knownvalue.ListExact([]knownvalue.Check{
			knownvalue.ObjectExact(map[string]knownvalue.Check{
				"alert_urgency_id": knownvalue.NotNull(),
				"type":             knownvalue.StringExact("escalation_policy"),
				"id":               knownvalue.NotNull(),
			}),
		})),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLiveCallRouterConfig(name, phoneNumber),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
				),
			},
			{
				Config: testAccResourceLiveCallRouterConfig(name+"-updated", phoneNumber),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
				),
			},
		},
	})
}

func testAccResourceLiveCallRouterConfig(name string, phoneNumber string) string {
	return testAccResourceEscalationPolicyConfig(name+"-live-call-router", name+"-live-call-router-description", "") + fmt.Sprintf(`
data "rootly_alert_urgency" "high" {
  name = "High"
}

resource "rootly_live_call_router" "test" {
	name = "%[1]s"

	country_code = "US"
	phone_type   = "local"
	phone_number = "%[2]s"

	voicemail_greeting = "Thank you for calling Rootly."

	paging_targets {
		alert_urgency_id = data.rootly_alert_urgency.high.id

		type = "escalation_policy"
		id   = rootly_escalation_policy.test.id
	}
}
`, name, phoneNumber)
}
