package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowCustomFieldSelection(t *testing.T) {
	workflowName := acctest.RandomWithPrefix("tf-workflow-custom-field")
	field1Name := acctest.RandomWithPrefix("tf-custom-field-1")
	field2Name := acctest.RandomWithPrefix("tf-custom-field-2")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowCustomFieldSelection(workflowName, field1Name, field2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_custom_field_selection.test1", "incident_condition", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_custom_field_selection.test2", "incident_condition", "ANY"),
				),
			},
		},
	})
}

func testAccResourceWorkflowCustomFieldSelection(workflowName, field1Name, field2Name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "test" {
  name = "%s"
}

resource "rootly_custom_field" "test1" {
	label = "%s"
	kind = "select"
	shown = ["incident_form", "incident_slack_form"]
	required = []
}

resource "rootly_custom_field" "test2" {
	label = "%s"
	kind = "text"
	shown = ["incident_form", "incident_slack_form"]
	required = []
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
`, workflowName, field1Name, field2Name)
}
