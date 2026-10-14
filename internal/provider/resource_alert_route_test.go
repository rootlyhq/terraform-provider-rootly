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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
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

// TER-227: removing `property_field_name` from a condition must clear it.
// Rules are recreated on every apply, so a stale value here means the provider
// is re-sending the removed attribute.
func TestAccResourceAlertRoute_RemovePropertyFieldName(t *testing.T) {
	resName := "rootly_alert_route.test"
	rName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConditionConfig(rName, `
					property_field_type   = "payload"
					property_field_name   = "$.team"
					property_field_values = ["infra"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("property_field_type"), knownvalue.StringExact("payload")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("property_field_name"), knownvalue.StringExact("$.team")),
				},
				Check: testAccCheckAlertRouteFirstCondition(resName, func(c alertRouteCondition) error {
					if c.PropertyFieldName != "$.team" {
						return fmt.Errorf("expected API property_field_name %q, got %q", "$.team", c.PropertyFieldName)
					}
					return nil
				}),
			},
			// Switch to the alert_field form and drop property_field_name entirely.
			{
				Config: testAccResourceAlertRouteConditionConfig(rName, `
					property_field_type   = "alert_field"
					conditionable_type    = "AlertField"
					conditionable_id      = data.rootly_alert_field.test.id
					property_field_values = ["infra"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("property_field_type"), knownvalue.StringExact("alert_field")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("conditionable_type"), knownvalue.StringExact("AlertField")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("property_field_name"), knownvalue.StringExact("")),
				},
				Check: testAccCheckAlertRouteFirstCondition(resName, func(c alertRouteCondition) error {
					if c.PropertyFieldName != "" {
						return fmt.Errorf("expected API property_field_name to be cleared, got %q", c.PropertyFieldName)
					}
					return nil
				}),
			},
		},
	})
}

// Same shape as TER-227 in the opposite direction: switching from the
// alert_field form back to payload must clear conditionable_type/conditionable_id.
func TestAccResourceAlertRoute_RemoveConditionable(t *testing.T) {
	resName := "rootly_alert_route.test"
	rName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConditionConfig(rName, `
					property_field_type   = "alert_field"
					conditionable_type    = "AlertField"
					conditionable_id      = data.rootly_alert_field.test.id
					property_field_values = ["infra"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("conditionable_type"), knownvalue.StringExact("AlertField")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("conditionable_id"), knownvalue.NotNull()),
				},
			},
			// Switch to the payload form and drop conditionable_type/conditionable_id.
			{
				Config: testAccResourceAlertRouteConditionConfig(rName, `
					property_field_type   = "payload"
					property_field_name   = "$.team"
					property_field_values = ["infra"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("property_field_type"), knownvalue.StringExact("payload")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("conditionable_type"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resName, testAccAlertRouteFirstConditionPath("conditionable_id"), knownvalue.StringExact("")),
				},
				Check: testAccCheckAlertRouteFirstCondition(resName, func(c alertRouteCondition) error {
					if c.ConditionableType != nil {
						return fmt.Errorf("expected API conditionable_type to be cleared, got %q", *c.ConditionableType)
					}
					if c.ConditionableId != nil {
						return fmt.Errorf("expected API conditionable_id to be cleared, got %q", c.ConditionableId.String())
					}
					return nil
				}),
			},
		},
	})
}

func testAccAlertRouteFirstConditionPath(attr string) tfjsonpath.Path {
	return tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("condition_groups").AtSliceIndex(0).AtMapKey("conditions").AtSliceIndex(0).AtMapKey(attr)
}

type alertRouteCondition struct {
	PropertyFieldName string
	ConditionableType *rootly.AlertRouteRulesConditionGroupsConditionsConditionableType
	ConditionableId   fmt.Stringer
}

// testAccCheckAlertRouteFirstCondition reads the alert route straight from the
// API, bypassing provider state, and checks the first rule's first condition.
func testAccCheckAlertRouteFirstCondition(resName string, check func(alertRouteCondition) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("not found: %s", resName)
		}

		httpResp, err := acctest.SharedClient.GetAlertRouteWithResponse(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting alert route, got error: %s", err)
		} else if httpResp.StatusCode() != http.StatusOK {
			return fmt.Errorf("Error getting alert route, got status code: %d", httpResp.StatusCode())
		} else if httpResp.ApplicationvndApiJSON200 == nil {
			return fmt.Errorf("Error getting alert route, got empty response")
		}

		rules := httpResp.ApplicationvndApiJSON200.Data.Attributes.Rules
		if rules == nil || len(*rules) == 0 || len((*rules)[0].ConditionGroups) == 0 || len((*rules)[0].ConditionGroups[0].Conditions) == 0 {
			return fmt.Errorf("expected alert route %s to have at least one condition", rs.Primary.ID)
		}

		c := (*rules)[0].ConditionGroups[0].Conditions[0]
		cond := alertRouteCondition{
			PropertyFieldName: c.PropertyFieldName,
			ConditionableType: c.ConditionableType,
		}
		if c.ConditionableId != nil {
			cond.ConditionableId = c.ConditionableId
		}
		return check(cond)
	}
}

func testAccResourceAlertRouteConditionConfig(name, condition string) string {
	return testAccResourceTeamConfig(name, "") + testAccResourceAlertUrgencyConfig(name, name+" description") + fmt.Sprintf(`
# Use a built-in alert field: route rule conditions are removed asynchronously,
# so deleting a custom alert field right after the route fails with a 422.
data "rootly_alert_field" "test" {
	kind = "title"
}

resource "rootly_escalation_policy" "test" {
	name      = "%[1]s"
	group_ids = [rootly_team.test.id]
}

resource "rootly_alerts_source" "test" {
	depends_on = [rootly_alert_urgency.test]

	name            = "%[1]s"
	source_type     = "generic_webhook"
	owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
	name              = "%[1]s"
	alerts_source_ids = [rootly_alerts_source.test.id]

	rules {
		name = "%[1]s-rule"

		destinations {
			target_type = "EscalationPolicy"
			target_id   = rootly_escalation_policy.test.id
		}

		condition_groups {
			conditions {
				property_field_condition_type = "is_one_of"
				%[2]s
			}
		}
	}
}
`, name, condition)
}
