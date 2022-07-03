package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendSlackMessage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendSlackMessage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskSendSlackMessageUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskSendSlackMessage = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_slack_message" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		text = "test"
	}
}
`

const testAccResourceWorkflowTaskSendSlackMessageUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_slack_message" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		text = "test"
	}
}
`
