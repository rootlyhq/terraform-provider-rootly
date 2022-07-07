package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateJiraIssue(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateJiraIssue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-workflow"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskCreateJiraIssueUpdate,
			},
		},
	})
}

const testAccResourceWorkflowTaskCreateJiraIssue = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_jira_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		project_key = "test"
title = "test"
issue_type = {
					id = "foo"
					name = "bar"
				}
	}
}
`

const testAccResourceWorkflowTaskCreateJiraIssueUpdate = `
resource "rootly_workflow_incident" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_create_jira_issue" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		project_key = "test"
title = "test"
issue_type = {
					id = "foo"
					name = "bar"
				}
	}
}
`