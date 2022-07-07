package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateLinearIssue(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateLinearIssue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateLinearIssueUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateLinearIssue = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_linear_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		title = "test"
team = {
					id = "foo"
					name = "bar"
				}
state = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskCreateLinearIssueUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_linear_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		title = "test"
team = {
					id = "foo"
					name = "bar"
				}
state = {
					id = "foo"
					name = "bar"
				}
	}
}
`