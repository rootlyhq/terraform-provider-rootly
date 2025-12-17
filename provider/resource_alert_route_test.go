package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoute(t *testing.T) {
	resName := "rootly_alert_route.test"
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	routeName := acctest.RandomWithPrefix("tf-alert-route")
	updatedRouteName := acctest.RandomWithPrefix("tf-alert-route-updated")
	ruleName := acctest.RandomWithPrefix("tf-rule")
	updatedRuleName := acctest.RandomWithPrefix("tf-rule-updated")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName, ruleName, "is_one_of", "$.severity", "critical"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", ruleName),
					resource.TestCheckResourceAttr(resName, "rules.0.position", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.fallback_rule", "false"),
				),
			},
			{
				Config: testAccResourceAlertRouteConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, updatedRouteName, updatedRuleName, "contains", "$.title", "error"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", updatedRouteName),
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

func TestAccResourceAlertRouteWithMultipleTeams(t *testing.T) {
	resName := "rootly_alert_route.multi_team"
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	primaryTeamName := acctest.RandomWithPrefix("tf-primary-team")
	secondaryTeamName := acctest.RandomWithPrefix("tf-secondary-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	routeName := acctest.RandomWithPrefix("tf-multi-team-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteWithMultipleTeamsConfig(alertUrgencyName, primaryTeamName, secondaryTeamName, alertsSourceName, routeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
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

func TestAccResourceAlertRouteDisabled(t *testing.T) {
	resName := "rootly_alert_route.disabled"
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	routeName := acctest.RandomWithPrefix("tf-disabled-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteDisabledConfig(alertUrgencyName, teamName, alertsSourceName, routeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
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
	resName := "rootly_alert_route.multi_rules"
	criticalEPName := acctest.RandomWithPrefix("tf-critical-escalation")
	warningEPName := acctest.RandomWithPrefix("tf-warning-escalation")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	routeName := acctest.RandomWithPrefix("tf-multi-rules-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteWithMultipleRulesConfig(criticalEPName, warningEPName, alertUrgencyName, teamName, alertsSourceName, routeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttr(resName, "rules.#", "3"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", "Critical Rule"),
					resource.TestCheckResourceAttr(resName, "rules.1.name", "Warning Rule"),
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
	resName := "rootly_alert_route.rules_update"
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	routeName := acctest.RandomWithPrefix("tf-rules-update-route")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteRulesUpdateBeforeConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
					resource.TestCheckResourceAttr(resName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "rules.0.name", "Initial Rule"),
				),
			},
			{
				Config: testAccResourceAlertRouteRulesUpdateAfterConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", routeName),
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

// Helper function to generate the basic alert route config with a single rule
func testAccResourceAlertRouteConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName, ruleName, conditionType, conditionField, conditionValue string) string {
	var conditionBlock string
	if conditionType == "is_one_of" {
		conditionBlock = fmt.Sprintf(`
        property_field_condition_type = "%s"
        property_field_name = "%s"
        property_field_type = "payload"
        property_field_values = ["%s"]`, conditionType, conditionField, conditionValue)
	} else {
		conditionBlock = fmt.Sprintf(`
        property_field_condition_type = "%s"
        property_field_name = "%s"
        property_field_type = "payload"
        property_field_value = "%s"`, conditionType, conditionField, conditionValue)
	}

	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
  description = "Test team for alert routing"
}

resource "rootly_escalation_policy" "production_ep" {
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

  rules {
    name = "%s"
    position = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.production_ep.id
    }

    condition_groups {
      position = 1

      conditions {
        %s
      }
    }
  }
}
`, teamName, escalationPolicyName, alertUrgencyName, alertsSourceName, routeName, ruleName, conditionBlock)
}

func testAccResourceAlertRouteWithMultipleTeamsConfig(alertUrgencyName, primaryTeamName, secondaryTeamName, alertsSourceName, routeName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test_primary" {
  name = "%s"
  description = "Primary team for alerts"
}

resource "rootly_team" "test_secondary" {
  name = "%s"
  description = "Secondary team for alerts"
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
  owner_group_ids = [rootly_team.test_primary.id, rootly_team.test_secondary.id]

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

resource "rootly_alert_route" "multi_team" {
  name = "%s"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test_primary.id, rootly_team.test_secondary.id]
}
`, primaryTeamName, secondaryTeamName, alertUrgencyName, alertsSourceName, routeName)
}

func testAccResourceAlertRouteDisabledConfig(alertUrgencyName, teamName, alertsSourceName, routeName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
  description = "Test team for alert routing"
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

resource "rootly_alert_route" "disabled" {
  name = "%s"
  enabled = false
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
}
`, teamName, alertUrgencyName, alertsSourceName, routeName)
}

func testAccResourceAlertRouteWithMultipleRulesConfig(criticalEPName, warningEPName, alertUrgencyName, teamName, alertsSourceName, routeName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
  description = "Test team for alert routing"
}

resource "rootly_escalation_policy" "critical_ep" {
  name      = "%s"
  group_ids = [rootly_team.test.id]
}

resource "rootly_escalation_policy" "warning_ep" {
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

resource "rootly_alert_route" "multi_rules" {
  name = "%s"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]

  rules {
    name = "Critical Rule"
    position = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.critical_ep.id
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
      target_id = rootly_escalation_policy.warning_ep.id
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
      target_id = rootly_team.test.id
    }
  }
}
`, teamName, criticalEPName, warningEPName, alertUrgencyName, alertsSourceName, routeName)
}

func testAccResourceAlertRouteRulesUpdateBeforeConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
  description = "Test team for alert routing"
}

resource "rootly_escalation_policy" "production_ep" {
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

resource "rootly_alert_route" "rules_update" {
  name = "%s"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]

  rules {
    name = "Initial Rule"
    position = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.production_ep.id
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
}
`, teamName, escalationPolicyName, alertUrgencyName, alertsSourceName, routeName)
}

func testAccResourceAlertRouteRulesUpdateAfterConfig(escalationPolicyName, alertUrgencyName, teamName, alertsSourceName, routeName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
  description = "Test team for alert routing"
}

resource "rootly_escalation_policy" "production_ep" {
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

resource "rootly_alert_route" "rules_update" {
  name = "%s"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]

  rules {
    name = "Updated Rule 1"
    position = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.production_ep.id
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
      target_id = rootly_team.test.id
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
}
`, teamName, escalationPolicyName, alertUrgencyName, alertsSourceName, routeName)
}
