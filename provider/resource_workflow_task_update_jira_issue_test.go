package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateJiraIssue(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-jira")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateJiraIssueConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_jira_issue.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_jira_issue.foo", "task_params.0.task_type", "update_jira_issue"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_jira_issue.foo", "task_params.0.issue_id", "PROJ-123"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_jira_issue.foo", "task_params.0.project_key", "PROJ"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskUpdateJiraIssueConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_jira_issue" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    issue_id    = "PROJ-123"
    project_key = "PROJ"
  }
}
`, name)
}
