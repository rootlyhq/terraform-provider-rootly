package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskArchiveSlackChannels(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskArchiveSlackChannels,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskArchiveSlackChannelsUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskArchiveSlackChannels = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_archive_slack_channels" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		channels {
						id = "foo"
						name = "bar"
					}
	}
}
`

const testAccResourceWorkflowTaskArchiveSlackChannelsUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_archive_slack_channels" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		channels {
						id = "foo"
						name = "bar"
					}
	}
}
`