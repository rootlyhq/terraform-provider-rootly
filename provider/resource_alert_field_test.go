package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertField(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertFieldCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", "Test Alert Field"),
				),
			},
			{
				Config: testAccResourceAlertFieldUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", "Updated Alert Field"),
				),
			},
		},
	})
}

const testAccResourceAlertFieldCreate = `
resource "rootly_alert_field" "test" {
  name = "Test Alert Field"
}
`

const testAccResourceAlertFieldUpdate = `
resource "rootly_alert_field" "test" {
  name = "Updated Alert Field"
}
`
