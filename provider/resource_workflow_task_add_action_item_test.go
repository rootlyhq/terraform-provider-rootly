package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAddActionItem(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAddActionItem,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskAddActionItemUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskAddActionItem = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_add_action_item" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		summary = "test"
status = "open"
priority = "low"
	}
}
`

const testAccResourceWorkflowTaskAddActionItemUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_add_action_item" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		summary = "test"
status = "open"
priority = "low"
	}
}
`