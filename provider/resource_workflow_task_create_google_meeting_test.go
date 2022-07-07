package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateGoogleMeeting(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateGoogleMeeting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateGoogleMeetingUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateGoogleMeeting = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_google_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		summary = "test"
description = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateGoogleMeetingUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_google_meeting" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		summary = "test"
description = "test"
	}
}
`