package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowFormFieldCondition(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-ffc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowFormFieldConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_form_field_condition.test1", "incident_condition", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_form_field_condition.test2", "incident_condition", "ANY"),
				),
			},
		},
	})
}

func testAccResourceWorkflowFormFieldConditionConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "test" {
  name = "%s-workflow"
}

resource "rootly_form_field" "test1" {
  	name = "%s-ff-1"
	kind = "custom"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field" "test2" {
  	name = "%s-ff-2"
	kind = "custom"
	input_kind = "text"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field_option" "test" {
	form_field_id = rootly_form_field.test1.id
	value = "test"
}

resource "rootly_workflow_form_field_condition" "test1" {
	workflow_id = rootly_workflow_incident.test.id
	form_field_id = rootly_form_field.test1.id
	selected_option_ids = [rootly_form_field_option.test.id]
}

resource "rootly_workflow_form_field_condition" "test2" {
	workflow_id = rootly_workflow_incident.test.id
	form_field_id = rootly_form_field.test2.id
	values = ["test"]
}
`, rName, rName, rName)
}
