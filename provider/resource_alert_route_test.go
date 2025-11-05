package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoute(t *testing.T) {
  resName := "rootly_alert_route.test"

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
					resource.TestCheckResourceAttr(resName, "name", "Test Alert Route"),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
      {
				ResourceName: resName,
				ImportState: true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRouteWithMultipleTeams(t *testing.T) {
  resName := "rootly_alert_route.multi_team"

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
					resource.TestCheckResourceAttr(resName, "name", "Multi Team Alert Route"),
					resource.TestCheckResourceAttr(resName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
      {
				ResourceName: resName,
				ImportState: true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertRouteDisabled(t *testing.T) {
  resName := "rootly_alert_route.disabled"

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
					resource.TestCheckResourceAttr(resName, "name", "Disabled Alert Route"),
					resource.TestCheckResourceAttr(resName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
      {
				ResourceName: resName,
				ImportState: true,
				ImportStateVerify: true,
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
  name = "Test Alerts Source"
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
  name = "Test Alerts Source"
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
  name = "Test Alerts Source"
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
