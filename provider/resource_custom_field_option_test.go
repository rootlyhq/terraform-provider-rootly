package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomFieldOption(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomFieldOption,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field.parent", "label", "mycustomfieldParent"),
					resource.TestCheckResourceAttr("rootly_custom_field_option.foo", "value", "myoption"),
				),
			},
			{
				Config: testAccResourceCustomFieldOptionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field_option.foo", "value", "myoption2"),
				),
			},
		},
	})
}

const testAccResourceCustomFieldOption = `
resource "rootly_custom_field" "parent" {
  label = "mycustomfieldParent"
	kind = "select"
	shown = ["incident_form", "incident_slack_form"]
	required = []
}

resource "rootly_custom_field_option" "foo" {
	custom_field_id = rootly_custom_field.parent.id
  value = "myoption"
}
`

const testAccResourceCustomFieldOptionUpdate = `
resource "rootly_custom_field" "parent" {
  label = "mycustomfieldParent"
	kind = "select"
	shown = ["incident_form", "incident_slack_form"]
	required = []
}

resource "rootly_custom_field_option" "foo" {
	custom_field_id = rootly_custom_field.parent.id
  value       = "myoption2"
}
`
