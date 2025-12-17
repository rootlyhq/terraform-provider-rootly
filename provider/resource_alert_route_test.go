package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoute(t *testing.T) {
	resName := "rootly_alert_route.test"
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")
	escalationPolicy1Name := acctest.RandomWithPrefix("tf-escalation-policy-1")
	escalationPolicy2Name := acctest.RandomWithPrefix("tf-escalation-policy-2")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          owning_team_ids = [rootly_team.test_1.id]

          rules {
            name = "High Priority Rule"
            position = 1
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "is_one_of"
                property_field_name = "$.severity"
                property_field_type = "payload"
                property_field_values = ["critical"]
              }
            }
          }
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", ruleName),
					resource.TestCheckResourceAttr(resName, "rules.0.position", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.fallback_rule", "false"),
				),
			},
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName+"-updated", `
          owning_team_ids = [rootly_team.test_1.id]

          rules {
            name = "Updated High Priority Rule"
            position = 1
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "contains"
                property_field_name = "$.title"
                property_field_type = "payload"
                property_field_value = "error"
              }
            }
          }
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName+"-updated"),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", updatedRuleName),
					resource.TestCheckResourceAttr(resName, "rules.0.position", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.fallback_rule", "false"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRoute_Validation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Empty rule, can validate
			{
				Config: `
          resource "rootly_alert_route" "test" {
            name              = "Test Alert Route"
            alerts_source_ids = ["alert-source-id"]

            rules {
              name = "Empty rule"
            }
          }
        `,
				ExpectError: regexp.MustCompile("Either `fallback_rule` or `condition_groups` must be specified"),
			},
			// Fallback rule, can validate
			{
				Config: `
          resource "rootly_alert_route" "test" {
            name              = "Test Alert Route"
            alerts_source_ids = ["alert-source-id"]

            rules {
              name = "Empty rule"
              fallback_rule = true
              condition_groups {
                position = 1
              }
            }
          }
        `,
				ExpectError: regexp.MustCompile("`fallback_rule` and `condition_groups` cannot be specified at the same time"),
			},
			// Unknown rules, cannot validate
			{
				Config: `
          resource "terraform_data" "rules" {
            input = [
              {
                name = "Empty rule"
                condition_groups = [
                  {
                    position = 1
                  }
                ]
              },
              {
                name = "Fallback rule"
                fallback_rule = true
              }
            ]
          }

          resource "rootly_alert_route" "test" {
            name              = "Test Alert Route"
            alerts_source_ids = ["alert-source-id"]

            dynamic "rules" {
              for_each = terraform_data.rules.output
              content {
                name = rules.value.name
                fallback_rule = rules.value.fallback_rule
                dynamic "condition_groups" {
                  for_each = rules.value.condition_groups
                  content {
                    position = condition_groups.value.position
                  }
                }
              }
            }
          }
        `,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Unknown rules, cannot validate
			{
				Config: `
          resource "terraform_data" "name" {
            input = "Empty rule"
          }

          resource "terraform_data" "fallback_rule" {
            input = true
          }

          resource "terraform_data" "condition_groups" {
            input = []
          }

          resource "rootly_alert_route" "test" {
            name              = "Test Alert Route"
            alerts_source_ids = ["alert-source-id"]

            rules {
              name = terraform_data.name.output
              fallback_rule = terraform_data.fallback_rule.output
              dynamic "condition_groups" {
                for_each = terraform_data.condition_groups.output
                content {
                }
              }
            }
          }
        `,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAlertRouteWithMultipleTeams(t *testing.T) {
	resName := "rootly_alert_route.test"
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")
	escalationPolicy1Name := acctest.RandomWithPrefix("tf-escalation-policy-1")
	escalationPolicy2Name := acctest.RandomWithPrefix("tf-escalation-policy-2")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          owning_team_ids = [rootly_team.test_1.id, rootly_team.test_2.id]
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "owning_team_ids.#", "2"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRouteDisabled(t *testing.T) {
	resName := "rootly_alert_route.test"
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")
	escalationPolicy1Name := acctest.RandomWithPrefix("tf-escalation-policy-1")
	escalationPolicy2Name := acctest.RandomWithPrefix("tf-escalation-policy-2")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          enabled = false
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRouteWithMultipleRules(t *testing.T) {
	resName := "rootly_alert_route.test"
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")
	escalationPolicy1Name := acctest.RandomWithPrefix("tf-escalation-policy-1")
	escalationPolicy2Name := acctest.RandomWithPrefix("tf-escalation-policy-2")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          rules {
            name = "Critical Rule"
            position = 1
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "is_one_of"
                property_field_name = "$.severity"
                property_field_type = "payload"
                property_field_values = ["critical"]
              }
            }
          }

          rules {
            name = "Warning Rule"
            position = 2
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_2.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "is_one_of"
                property_field_name = "$.severity"
                property_field_type = "payload"
                property_field_values = ["warning"]
              }
            }
          }

          rules {
            position = 3
            fallback_rule = true

            destinations {
              target_type = "Group"
              target_id = rootly_team.test_1.id
            }
          }
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "rules.#", "3"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", "Critical Rule"),
					resource.TestCheckResourceAttr(resName, "rules.1.name", "Warning Rule"),
					resource.TestCheckResourceAttr(resName, "rules.2.name", "Fallback Rule for "+alertRouteName),
					resource.TestCheckResourceAttr(resName, "rules.2.fallback_rule", "true"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRouteRulesUpdate(t *testing.T) {
	resName := "rootly_alert_route.test"
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")
	escalationPolicy1Name := acctest.RandomWithPrefix("tf-escalation-policy-1")
	escalationPolicy2Name := acctest.RandomWithPrefix("tf-escalation-policy-2")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertRouteName := acctest.RandomWithPrefix("tf-alert-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          owning_team_ids = [rootly_team.test_1.id]

          rules {
            name = "Initial Rule"
            position = 1
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "is_one_of"
                property_field_name = "$.severity"
                property_field_type = "payload"
                property_field_values = ["critical"]
              }
            }
          }
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", "Initial Rule"),
				),
			},
			{
				Config: testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, `
          owning_team_ids = [rootly_team.test_1.id]

          rules {
            name = "Updated Rule 1"
            position = 1
            fallback_rule = false

            destinations {
              target_type = "EscalationPolicy"
              target_id = rootly_escalation_policy.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "contains"
                property_field_name = "$.title"
                property_field_type = "payload"
                property_field_value = "error"
              }
            }
          }

          rules {
            name = "New Rule 2"
            position = 2
            fallback_rule = false

            destinations {
              target_type = "Group"
              target_id = rootly_team.test_1.id
            }

            condition_groups {
              position = 1

              conditions {
                property_field_condition_type = "is_one_of"
                property_field_name = "$.severity"
                property_field_type = "payload"
                property_field_values = ["warning"]
              }
            }
          }
        `),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertRouteName),
					resource.TestCheckResourceAttr(resName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", "Updated Rule 1"),
					resource.TestCheckResourceAttr(resName, "rules.1.name", "New Rule 2"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceAlertRouteConfig(team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, extra string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test_1" {
  name = "%[1]s"
  description = "%[1]s-description"
}

resource "rootly_team" "test_2" {
  name = "%[2]s"
  description = "%[2]s-description"
}

resource "rootly_escalation_policy" "test_1" {
  name      = "%[3]s"
  group_ids = [rootly_team.test_1.id, rootly_team.test_2.id]
}

resource "rootly_escalation_policy" "test_2" {
  name      = "%[4]s"
  group_ids = [rootly_team.test_1.id, rootly_team.test_2.id]
}

resource "rootly_alert_urgency" "test" {
  name = "%[5]s"
  description = "%[5]s-description"
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
  name = "%[6]s"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test_1.id, rootly_team.test_2.id]

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
  name = "%[7]s"
  alerts_source_ids = [rootly_alerts_source.test.id]

  %[8]s
}
`, team1Name, team2Name, escalationPolicy1Name, escalationPolicy2Name, alertUrgencyName, alertsSourceName, alertRouteName, extra)
}
