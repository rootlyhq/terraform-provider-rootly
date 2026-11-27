package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendSlackBlocks(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-wf-task")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendSlackBlocksConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", rName),
				),
			},
			{
				Config: testAccResourceWorkflowTaskSendSlackBlocksUpdateConfig(rName),
			},
		},
	})
}

func testAccResourceWorkflowTaskSendSlackBlocksConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  	name = "%s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_slack_blocks" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		blocks = "test"
channels {
						id = "foo"
						name = "bar"
					}
retry_wait_time = 1
	}
}
`, name)
}

func testAccResourceWorkflowTaskSendSlackBlocksUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  	name = "%s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_slack_blocks" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		blocks = "test"
channels {
						id = "foo"
						name = "bar"
					}
retry_wait_time = 1
	}
}
`, name)
}

// blocks written as a heredoc end with a newline, which the API trims when the
// value is not valid JSON (e.g. it contains Liquid tags). This must not cause a
// perpetual diff.
func TestAccResourceWorkflowTaskSendSlackBlocks_liquidHeredoc(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-wf-task")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// The test framework fails the step if the plan after apply is not empty.
				Config: fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
	name = "%s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_slack_blocks" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		channels {
			id   = "foo"
			name = "bar"
		}
		blocks = <<-EOT
			[
			  {"type": "section", "text": {"type": "mrkdwn", "text": "{{ incident.title }}"}}{%% if incident.severity %%},
			  {"type": "divider"}{%% endif %%}
			]
		EOT
	}
}
`, rName),
			},
		},
	})
}
