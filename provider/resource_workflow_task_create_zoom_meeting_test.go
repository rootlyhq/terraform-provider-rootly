package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateZoomMeeting(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateZoomMeeting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateZoomMeetingUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateZoomMeeting = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_zoom_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		topic = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateZoomMeetingUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_zoom_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		topic = "test"
	}
}
`