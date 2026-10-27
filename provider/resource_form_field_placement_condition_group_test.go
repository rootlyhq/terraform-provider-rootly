package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormFieldPlacementConditionGroup(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-condition-group")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormFieldPlacementConditionGroup(name, "or"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field_placement_condition_group.test", "group_type", "required"),
					resource.TestCheckResourceAttr("rootly_form_field_placement_condition_group.test", "group_operator", "or"),
					resource.TestCheckResourceAttrPair("rootly_form_field_placement_condition.option_a", "condition_group_id", "rootly_form_field_placement_condition_group.test", "id"),
					resource.TestCheckResourceAttrPair("rootly_form_field_placement_condition.option_b", "condition_group_id", "rootly_form_field_placement_condition_group.test", "id"),
				),
			},
			{
				Config: testAccResourceFormFieldPlacementConditionGroup(name, "and"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field_placement_condition_group.test", "group_operator", "and"),
					resource.TestCheckResourceAttrPair("rootly_form_field_placement_condition.option_b", "condition_group_id", "rootly_form_field_placement_condition_group.test", "id"),
				),
			},
		},
	})
}

func testAccResourceFormFieldPlacementConditionGroup(name, groupOperator string) string {
	return fmt.Sprintf(`
resource "rootly_form_set" "test" {
  name  = "%[1]s"
  forms = ["web_new_incident_form"]
}

resource "rootly_form_field" "source" {
  name       = "%[1]s-source"
  input_kind = "select"
  shown      = ["web_new_incident_form"]
  required   = []
}

resource "rootly_form_field_option" "a" {
  form_field_id = rootly_form_field.source.id
  value         = "a"
}

resource "rootly_form_field_option" "b" {
  form_field_id = rootly_form_field.source.id
  value         = "b"
}

resource "rootly_form_field" "target" {
  name       = "%[1]s-target"
  input_kind = "text"
  shown      = ["web_new_incident_form"]
  required   = []
}

resource "rootly_form_field_placement" "test" {
  form_field_id = rootly_form_field.target.id
  form_set_id   = rootly_form_set.test.id
  form          = "web_new_incident_form"
}

resource "rootly_form_field_placement_condition_group" "test" {
  form_field_placement_id = rootly_form_field_placement.test.id
  group_type              = "required"
  group_operator          = "%[2]s"
}

resource "rootly_form_field_placement_condition" "option_a" {
  form_field_placement_id = rootly_form_field_placement.test.id
  condition_group_id      = rootly_form_field_placement_condition_group.test.id
  conditioned             = "required"
  form_field_id           = rootly_form_field.source.id
  comparison              = "equal"
  values                  = [rootly_form_field_option.a.id]
}

resource "rootly_form_field_placement_condition" "option_b" {
  form_field_placement_id = rootly_form_field_placement.test.id
  condition_group_id      = rootly_form_field_placement_condition_group.test.id
  conditioned             = "required"
  form_field_id           = rootly_form_field.source.id
  comparison              = "equal"
  values                  = [rootly_form_field_option.b.id]
}
`, name, groupOperator)
}
