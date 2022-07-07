package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskInviteToSlackChannelVictorOps(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskInviteToSlackChannelVictorOps,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskInviteToSlackChannelVictorOpsUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskInviteToSlackChannelVictorOps = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_invite_to_slack_channel_victor_ops" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		schedule = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskInviteToSlackChannelVictorOpsUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_invite_to_slack_channel_victor_ops" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		schedule = {
					id = "foo"
					name = "bar"
				}
	}
}
`