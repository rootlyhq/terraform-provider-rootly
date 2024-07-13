package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormFieldOption(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormFieldOption,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field.parent", "name", "myformfieldParent"),
					resource.TestCheckResourceAttr("rootly_form_field_option.foo", "value", "myoption"),
				),
			},
			{
				Config: testAccResourceFormFieldOptionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field_option.foo", "value", "myoption2"),
				),
			},
		},
	})
}

const testAccResourceFormFieldOption = `
resource "rootly_form_field" "parent" {
  name = "myformfieldParent"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field_option" "foo" {
	form_field_id = rootly_form_field.parent.id
  value = "myoption"
}
`

const testAccResourceFormFieldOptionUpdate = `
resource "rootly_form_field" "parent" {
  name = "myformfieldParent"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field_option" "foo" {
	form_field_id = rootly_form_field.parent.id
  value       = "myoption2"
}
`
