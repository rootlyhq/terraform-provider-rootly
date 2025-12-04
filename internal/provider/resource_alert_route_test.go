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
	resource.AddTestSweepers("rootly_alert_route", &resource.Sweeper{
		Name: "rootly_alert_route",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListAlertRoutesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListAlertRoutesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting alert routes, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting alert routes, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting alert routes, got empty response")
				}

				for _, alertRoute := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(alertRoute.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteAlertRouteWithResponse(ctx, alertRoute.Id)
						if err != nil {
							return fmt.Errorf("Error deleting alert route: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting alert route, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted alert route %s", alertRoute.Attributes.Name)
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

func TestAccResourceAlertRouteV2(t *testing.T) {
	resName := "rootly_alert_route.test"
	teamName := acctest.RandomWithPrefix("tf-team")
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(alertRouteName)),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccResourceAlertRoute(teamName, escalationPolicyName, alertUrgencyName, alertsSourceName, alertRouteName),
				ConfigStateChecks: append(configStateChecks),
			},
			{
				Config:            testAccResourceAlertRoute(teamName, escalationPolicyName, alertUrgencyName, alertsSourceName, alertRouteName),
				ConfigStateChecks: append(configStateChecks),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceAlertRoute(teamName string, escalationPolicyName string, alertUrgencyName string, alertsSourceName string, alertRouteName string) string {
	return testAccResourceTeam(teamName) + fmt.Sprintf(`
resource "rootly_escalation_policy" "test" {
  name      = "%s"
  group_ids = [rootly_team.test.id]
}

resource "rootly_alert_urgency" "test" {
  name = "%s"
  description = "Test urgency for alerts"
  position = 1
}

data "rootly_alert_field" "title_field" {
  kind = "title"
}

data "rootly_alert_field" "description_field" {
  kind = "description"
}

data "rootly_alert_field" "source_link_field" {
  kind = "external_url"
}

resource "rootly_alerts_source" "test" {
  name = "%s"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.title_field.id
  }

  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.description_field.id
  }

  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.source_link_field.id
  }

  alert_source_urgency_rules_attributes {
    alert_urgency_id = rootly_alert_urgency.test.id
    json_path = "severity"
    operator = "is"
    value = "critical"
  }

  sourceable_attributes {
    auto_resolve = false
    resolve_state = "$.resolved"
  }
}

resource "rootly_alert_route" "test" {
  name = "%s"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]

  rules = [
		{
			name = "High Priority Rule"
			position = 1
			fallback_rule = false

			destinations = [
				{
					target_type = "EscalationPolicy"
					target_id = rootly_escalation_policy.test.id
				}
			]

			condition_groups = [
				{
					position = 1

					conditions = [
						{
							property_field_condition_type = "is_one_of"
							property_field_name = "$.severity"
							property_field_type = "payload"
							property_field_values = ["critical"]
						}
					]
				}
			]
		}
	]
}
`, escalationPolicyName, alertUrgencyName, alertsSourceName, alertRouteName)
}
