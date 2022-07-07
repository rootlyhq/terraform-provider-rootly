package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateJiraSubtask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateJiraSubtask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateJiraSubtaskUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateJiraSubtask = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_jira_subtask" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		project_key = "test"
parent_issue_id = "test"
title = "test"
subtask_issue_type = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskCreateJiraSubtaskUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_jira_subtask" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		project_key = "test"
parent_issue_id = "test"
title = "test"
subtask_issue_type = {
					id = "foo"
					name = "bar"
				}
	}
}
`