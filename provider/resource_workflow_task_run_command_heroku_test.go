package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskRunCommandHeroku(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskRunCommandHeroku,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskRunCommandHerokuUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskRunCommandHeroku = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_run_command_heroku" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		command = "test"
app_name = "test"
size = "standard-1X"
	}
}
`

const testAccResourceWorkflowTaskRunCommandHerokuUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_run_command_heroku" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		command = "test"
app_name = "test"
size = "standard-1X"
	}
}
`