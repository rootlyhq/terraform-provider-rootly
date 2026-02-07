package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormFieldOption(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-form-field")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormFieldOption(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field.parent", "name", randomName),
					resource.TestCheckResourceAttr("rootly_form_field_option.foo", "value", "myoption"),
				),
			},
			{
				Config: testAccResourceFormFieldOptionUpdate(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field_option.foo", "value", "myoption2"),
				),
			},
		},
	})
}

func testAccResourceFormFieldOption(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_field" "parent" {
  name = "%s"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field_option" "foo" {
	form_field_id = rootly_form_field.parent.id
  value = "myoption"
}
`, name)
}

func testAccResourceFormFieldOptionUpdate(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_field" "parent" {
  name = "%s"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field_option" "foo" {
	form_field_id = rootly_form_field.parent.id
  value       = "myoption2"
}
`, name)
}
