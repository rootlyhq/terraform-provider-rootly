package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomField(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
			time.Sleep(1 * time.Second)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "label", "mycustomfield"),
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceCustomFieldUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "label", "mycustomfield2"),
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_custom_field.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceCustomField = `
resource "rootly_custom_field" "foo" {
  label = "mycustomfield"
	kind = "select"
	shown = ["incident_form"]
	required = ["incident_form"]
}
`

const testAccResourceCustomFieldUpdate = `
resource "rootly_custom_field" "foo" {
  label       = "mycustomfield2"
	kind = "select"
  description = "test description"
  enabled     = false
	shown = ["incident_form", "incident_post_mortem"]
	required = ["incident_form"]
}
`
