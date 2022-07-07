package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendEmail(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendEmail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskSendEmailUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskSendEmail = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_email" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		to = ["foo"]
subject = "test"
body = "test"
	}
}
`

const testAccResourceWorkflowTaskSendEmailUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_email" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		to = ["foo"]
subject = "test"
body = "test"
	}
}
`