package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateServiceNowIncident(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateServiceNowIncident,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateServiceNowIncidentUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateServiceNowIncident = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_service_now_incident" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		incident_id = "test"
	}
}
`

const testAccResourceWorkflowTaskUpdateServiceNowIncidentUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_service_now_incident" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		incident_id = "test"
	}
}
`