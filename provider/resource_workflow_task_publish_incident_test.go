package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskPublishIncident(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskPublishIncident,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskPublishIncidentUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskPublishIncident = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_publish_incident" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		incident = {
					id = "foo"
					name = "bar"
				}
public_title = "test"
event = "test"
status = "investigating"
status_page_ids = ["foo"]
	}
}
`

const testAccResourceWorkflowTaskPublishIncidentUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_publish_incident" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		incident = {
					id = "foo"
					name = "bar"
				}
public_title = "test"
event = "test"
status = "investigating"
status_page_ids = ["foo"]
	}
}
`