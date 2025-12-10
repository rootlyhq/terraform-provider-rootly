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
	rn := "rootly_escalation_policy.test"
	teamName := acctest.RandomWithPrefix("tf-team")
	policyName := acctest.RandomWithPrefix("tf-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPolicyConfig(teamName, policyName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(policyName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("business_hours"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"days": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("M"),
							knownvalue.StringExact("T"),
							knownvalue.StringExact("W"),
							knownvalue.StringExact("R"),
							knownvalue.StringExact("F"),
						}),
						"start_time": knownvalue.StringExact("09:00"),
						"end_time":   knownvalue.StringExact("17:00"),
						"time_zone":  knownvalue.Null(),
					})),
				},
			},
			{
				Config: testAccResourceEscalationPolicyConfig(teamName, policyName+"-updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(policyName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("business_hours"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"days": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("M"),
							knownvalue.StringExact("T"),
							knownvalue.StringExact("W"),
							knownvalue.StringExact("R"),
							knownvalue.StringExact("F"),
						}),
						"start_time": knownvalue.StringExact("09:00"),
						"end_time":   knownvalue.StringExact("17:00"),
						"time_zone":  knownvalue.Null(),
					})),
				},
			},
		},
	})
}

func testAccResourceEscalationPolicyConfig(teamName, policyName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}

resource "rootly_escalation_policy" "test" {
	name      = "%s"
	group_ids = [rootly_team.test.id]

	business_hours = {
		days = ["M", "T", "W", "R", "F"]
		start_time = "09:00"
		end_time = "17:00"
		time_zone = null
	}
}
	`, teamName, policyName)
}
