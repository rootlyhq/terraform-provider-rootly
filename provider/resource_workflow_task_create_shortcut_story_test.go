package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateShortcutStory(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateShortcutStory,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateShortcutStoryUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateShortcutStory = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_shortcut_story" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		title = "test"
project = {
					id = "foo"
					name = "bar"
				}
archivation = {
					id = "foo"
					name = "bar"
				}
kind = "bug"
	}
}
`

const testAccResourceWorkflowTaskCreateShortcutStoryUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_shortcut_story" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		title = "test"
project = {
					id = "foo"
					name = "bar"
				}
archivation = {
					id = "foo"
					name = "bar"
				}
kind = "bug"
	}
}
`