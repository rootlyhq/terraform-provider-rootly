package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowFormFieldCondition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowFormFieldCondition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_form_field_condition.test1", "incident_condition", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_form_field_condition.test2", "incident_condition", "ANY"),
				),
			},
		},
	})
}

const testAccResourceWorkflowFormFieldCondition = `
resource "rootly_workflow_incident" "test" {
  name = "workflow-form-field-test"
}

resource "rootly_form_field" "test1" {
  	name = "form-field-test-1"
	kind = "custom"
	input_kind = "select"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_form_field" "test2" {
  	name = "form-field-test-2"
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
`

func TestAccResourceWorkflowFormFieldConditionWithIsNotCondition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowFormFieldConditionWithIsNotCondition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_form_field_condition.test_is_not", "incident_condition", "IS NOT"),
				),
			},
		},
	})
}

const testAccResourceWorkflowFormFieldConditionWithIsNotCondition = `
resource "rootly_workflow_incident" "test" {
  name = "workflow-form-field-test-is-not"
}

resource "rootly_form_field" "test" {
  	name = "form-field-test-is-not"
	kind = "custom"
	input_kind = "text"
	shown = ["web_new_incident_form", "slack_new_incident_form"]
	required = []
}

resource "rootly_workflow_form_field_condition" "test_is_not" {
	workflow_id = rootly_workflow_incident.test.id
	form_field_id = rootly_form_field.test.id
	values = ["test"]
	incident_condition = "IS NOT"
}
`
