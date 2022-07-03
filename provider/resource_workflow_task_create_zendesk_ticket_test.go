package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateZendeskTicket(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateZendeskTicket,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateZendeskTicketUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateZendeskTicket = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_zendesk_ticket" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		kind = "problem"
subject = "test"
	}
}
`

const testAccResourceWorkflowTaskCreateZendeskTicketUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_zendesk_ticket" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		kind = "problem"
subject = "test"
	}
}
`
