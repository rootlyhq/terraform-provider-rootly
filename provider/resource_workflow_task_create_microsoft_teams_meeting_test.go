package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateMicrosoftTeamsMeeting(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateMicrosoftTeamsMeeting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateMicrosoftTeamsMeetingUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateMicrosoftTeamsMeeting = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_microsoft_teams_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		name = "test"
subject = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateMicrosoftTeamsMeetingUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_microsoft_teams_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		name = "test"
subject = "test"
	}
}
`