package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateOpsgenieIncident(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateOpsgenieIncident,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.task_type", "update_opsgenie_incident"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.opsgenie_incident_id", "12345"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.priority", "P1"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateOpsgenieIncidentEmptyPriority,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.priority", ""),
				),
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateOpsgenieIncident = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow-for-opsgenie-task"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_opsgenie_incident" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    opsgenie_incident_id = "12345"
    priority             = "P1"
  }
}
`

const testAccResourceWorkflowTaskUpdateOpsgenieIncidentEmptyPriority = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow-for-opsgenie-task"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_opsgenie_incident" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    opsgenie_incident_id = "12345"
    priority             = ""
  }
}
`
