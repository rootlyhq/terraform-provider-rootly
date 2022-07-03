package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateAsanaTask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateAsanaTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateAsanaTaskUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateAsanaTask = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_asana_task" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		task_id = "test"
completion = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskUpdateAsanaTaskUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_asana_task" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		task_id = "test"
completion = {
					id = "foo"
					name = "bar"
				}
	}
}
`
