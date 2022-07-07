package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateLinearIssue(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateLinearIssue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateLinearIssueUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateLinearIssue = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_linear_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		issue_id = "test"
	}
}
`

const testAccResourceWorkflowTaskUpdateLinearIssueUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_linear_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		issue_id = "test"
	}
}
`