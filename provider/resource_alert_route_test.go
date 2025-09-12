package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Test Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.test", "id"),
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.test", "id"),
				),
			},
		},
	})
}

func TestAccResourceAlertRouteWithMultipleTeams(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteWithMultipleTeams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.multi_team", "name", "Multi Team Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.multi_team", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.multi_team", "id"),
				),
			},
		},
	})
}

func TestAccResourceAlertRouteDisabled(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteDisabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.disabled", "name", "Disabled Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.disabled", "enabled", "false"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.disabled", "id"),
				),
			},
		},
	})
}

const testAccResourceAlertRouteCreate = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alert_urgency" "test" {
  name = "Test Alert Urgency"
  description = "Test urgency for alerts"
  position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

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
  name = "Test Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
}
`

const testAccResourceAlertRouteUpdate = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alert_urgency" "test" {
  name = "Test Alert Urgency"
  description = "Test urgency for alerts"
  position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

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
  name = "Updated Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
}
`

const testAccResourceAlertRouteWithMultipleTeams = `
resource "rootly_team" "test_primary" {
  name = "Primary Team"
  description = "Primary team for alerts"
}

resource "rootly_team" "test_secondary" {
  name = "Secondary Team"
  description = "Secondary team for alerts"
}

resource "rootly_alert_urgency" "test" {
  name = "Test Alert Urgency"
  description = "Test urgency for alerts"
  position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test_primary.id, rootly_team.test_secondary.id]

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
  name = "Multi Team Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test_primary.id, rootly_team.test_secondary.id]
}
`

const testAccResourceAlertRouteDisabled = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alert_urgency" "test" {
  name = "Test Alert Urgency"
  description = "Test urgency for alerts"
  position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

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
  name = "Disabled Alert Route"
  enabled = false
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
}
`