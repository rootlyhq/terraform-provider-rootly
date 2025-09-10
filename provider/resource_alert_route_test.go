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
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceAlertRouteCreate = `
resource "rootly_team" "test" {
	name = "test-team"
	description = "Test team for alert route"
}

resource "rootly_alerts_source" "test" {
	name = "test-alerts-source"
	source_type = "generic_webhook"
	owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
	name = "Test Alert Route"
	enabled = true
	alerts_source_ids = [rootly_alerts_source.test.id]
	owner_group_ids = [rootly_team.test.id]
}
`

const testAccResourceAlertRouteUpdate = `
resource "rootly_team" "test" {
	name = "test-team"
	description = "Test team for alert route"
}

resource "rootly_alerts_source" "test" {
	name = "test-alerts-source"
	source_type = "generic_webhook"
	owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
	name = "Updated Alert Route"
	enabled = false
	alerts_source_ids = [rootly_alerts_source.test.id]
	owner_group_ids = [rootly_team.test.id]
}
`
