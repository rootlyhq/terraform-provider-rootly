package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateZendeskTicket(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateZendeskTicket,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateZendeskTicketUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateZendeskTicket = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_zendesk_ticket" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		ticket_id = "test"
	}
}
`

const testAccResourceWorkflowTaskUpdateZendeskTicketUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_zendesk_ticket" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		ticket_id = "test"
	}
}
`