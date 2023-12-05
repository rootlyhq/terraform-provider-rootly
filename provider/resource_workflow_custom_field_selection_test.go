package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowCustomFieldSelection(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowCustomFieldSelection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_custom_field_selection.test1", "incident_condition", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_custom_field_selection.test2", "incident_condition", "ANY"),
				),
			},
		},
	})
}

const testAccResourceWorkflowCustomFieldSelection = `
resource "rootly_workflow_incident" "test" {
  name = "workflow-custom-field-test"
}

resource "rootly_custom_field" "test1" {
  	label = "custom-field-test-1"
	kind = "select"
	shown = ["incident_form"]
	required = ["incident_form"]
}

resource "rootly_custom_field" "test2" {
  	label = "custom-field-test-2"
	kind = "text"
	shown = ["incident_form"]
	required = ["incident_form"]
}

resource "rootly_custom_field_option" "test" {
	custom_field_id = rootly_custom_field.test1.id
	value = "test"
}

resource "rootly_workflow_custom_field_selection" "test1" {
	workflow_id = rootly_workflow_incident.test.id
	custom_field_id = rootly_custom_field.test1.id
	selected_option_ids = [rootly_custom_field_option.test.id]
}

resource "rootly_workflow_custom_field_selection" "test2" {
	workflow_id = rootly_workflow_incident.test.id
	custom_field_id = rootly_custom_field.test2.id
	values = ["test"]
}
`
