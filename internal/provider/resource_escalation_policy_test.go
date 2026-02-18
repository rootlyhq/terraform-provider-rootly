package provider

import (
	"context"
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
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func init() {
	resource.AddTestSweepers("rootly_escalation_policy", &resource.Sweeper{
		Name: "rootly_escalation_policy",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListEscalationPoliciesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListEscalationPoliciesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting escalation policies, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting escalation policies, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting escalation policies, got empty response")
				}

				for _, escalationPolicy := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(escalationPolicy.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteEscalationPolicyWithResponse(ctx, escalationPolicy.Id)
						if err != nil {
							return fmt.Errorf("Error deleting escalation policy: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting escalation policy, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted escalation policy %s", escalationPolicy.Attributes.Name)
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

func TestAccResourceEscalationPolicy(t *testing.T) {
	resName := "rootly_escalation_policy.test"
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPolicyConfig(escalationPolicyName, "description", `
					business_hours {
						time_zone  = "America/New_York"
						start_time = "12:00 PM"
						end_time   = "13:00"
						days       = ["M", "T", "W"]
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(escalationPolicyName)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("description")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("time_zone"), knownvalue.StringExact("America/New_York")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("start_time"), knownvalue.StringExact("12:00")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("end_time"), knownvalue.StringExact("13:00")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("days"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("M"),
						knownvalue.StringExact("T"),
						knownvalue.StringExact("W"),
					})),
				),
			},
			{
				Config: testAccResourceEscalationPolicyConfig(escalationPolicyName+"-updated", "description-updated", `
					business_hours {
						time_zone  = "America/Los_Angeles"
						start_time = "09:00"
						end_time   = "10:00"
						days       = ["W", "R", "F"]
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(escalationPolicyName+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("description-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("time_zone"), knownvalue.StringExact("America/Los_Angeles")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("start_time"), knownvalue.StringExact("09:00")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("end_time"), knownvalue.StringExact("10:00")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("business_hours").AtSliceIndex(0).AtMapKey("days"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("W"),
						knownvalue.StringExact("R"),
						knownvalue.StringExact("F"),
					})),
				),
			},
		},
	})
}

func testAccResourceEscalationPolicyConfig(name, description, extra string) string {
	return fmt.Sprintf(`
resource "rootly_escalation_policy" "test" {
	name        = "%s"
	description = "%s"
	%s
}
`, name, description, extra)
}
