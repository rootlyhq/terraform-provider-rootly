package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateShortcutTask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateShortcutTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateShortcutTaskUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateShortcutTask = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_shortcut_task" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		task_id = "test"
parent_story_id = "test"
completion = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskUpdateShortcutTaskUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_shortcut_task" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		task_id = "test"
parent_story_id = "test"
completion = {
					id = "foo"
					name = "bar"
				}
	}
}
`