package provider

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func testAccWorkflowTaskSlackCanvas(t *testing.T, action string) {
	t.Helper()
	if os.Getenv("ROOTLY_TEST_SLACK_CANVAS") != "1" {
		t.Skip("Set ROOTLY_TEST_SLACK_CANVAS=1 for an organization with Slack Canvas workflow actions enabled")
	}
	rName := acctest.RandomWithPrefix("tf-prf449-" + action)
	address := "rootly_workflow_task_" + action + ".test"
	initial := testAccWorkflowTaskSlackCanvasConfig(action, rName, false, false)
	selected := testAccWorkflowTaskSlackCanvasConfig(action, rName, false, true)
	updated := testAccWorkflowTaskSlackCanvasConfig(action, rName, true, true)
	cleared := testAccWorkflowTaskSlackCanvasConfig(action, rName, true, false)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("ROOTLY_API_TOKEN") == "" {
				t.Fatal("ROOTLY_API_TOKEN is required")
			}
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(state *terraform.State) error {
			api, err := client.NewClient(canvasAcceptanceAPIURL(), os.Getenv("ROOTLY_API_TOKEN"), "canvas-acceptance-test")
			if err != nil {
				return err
			}
			for _, record := range state.RootModule().Resources {
				switch record.Type {
				case "rootly_workflow_task_" + action:
					_, err = api.GetWorkflowTask(record.Primary.ID)
				case "rootly_workflow_incident":
					_, err = api.GetWorkflow(record.Primary.ID)
				default:
					continue
				}
				if !errors.Is(err, client.NewNotFoundError("")) {
					return fmt.Errorf("Canvas workflow resource still exists or cleanup failed: %v", err)
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: initial, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("rootly_workflow_incident.test", "enabled", "false"),
				resource.TestCheckResourceAttr(address, "enabled", "false"),
				resource.TestCheckResourceAttr(address, "task_params.0.retry_count", "0"),
				resource.TestCheckResourceAttr(address, "task_params.0.retry_wait_time", "1"),
				testAccWorkflowTaskSlackCanvasActionCheck(action, address, rName, false),
				resource.TestCheckResourceAttr(address, "task_params.0.content", "# Incident {{ incident.title }}\nInitial report"),
				resource.TestCheckResourceAttr(address, "task_params.0.channel.0.workspace.#", "0"),
			)},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: initial, PlanOnly: true},
			{Config: selected, Check: resource.TestCheckResourceAttr(address, "task_params.0.channel.0.workspace.0.id", "T_CANVAS_EAST")},
			{Config: updated, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "task_params.0.content", "# Incident {{ incident.title }}\nUpdated report"),
				resource.TestCheckResourceAttr(address, "task_params.0.channel.0.workspace.0.id", "T_CANVAS_WEST"),
				resource.TestCheckResourceAttr(address, "task_params.0.channel.0.id", "C_CANVAS_WEST"),
				testAccWorkflowTaskSlackCanvasActionCheck(action, address, rName, true),
			)},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: updated, PlanOnly: true},
			{Config: cleared, Check: resource.TestCheckResourceAttr(address, "task_params.0.channel.0.workspace.#", "0")},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: cleared, PlanOnly: true},
		},
	})
}

func testAccWorkflowTaskSlackCanvasConfig(action, name string, updated, workspace bool) string {
	content, suffix := "Initial report", "EAST"
	trigger := "incident_updated"
	if action == "create_slack_canvas" {
		trigger = "slack_channel_created"
	}
	if updated {
		content, suffix = "Updated report", "WEST"
	}
	workspaceBlock := ""
	if workspace {
		workspaceBlock = fmt.Sprintf(`workspace {
        id = "T_CANVAS_%s"
        name = "Canvas %s"
      }`, suffix, suffix)
	}
	specific := ""
	if action == "create_slack_canvas" {
		title := name + " {{ incident.sequential_id }}"
		if updated {
			title = "Updated " + title
		}
		specific = fmt.Sprintf("title = %q", title)
	}
	if action == "update_slack_canvas" && updated {
		specific = `operation = "managed_sections"`
	}
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "test" {
  name = %q
  enabled = false
  trigger_params { triggers = [%q] }
}
resource "rootly_workflow_task_%s" "test" {
  workflow_id = rootly_workflow_incident.test.id
  name = %q
  enabled = false
  position = 1
  task_params {
    %s
    channel {
      id = "C_CANVAS_%s"
      name = "canvas-test"
      %s
    }
    content = "# Incident {{ incident.title }}\n%s"
  }
}
`, name, trigger, action, name, specific, suffix, workspaceBlock, content)
}

func testAccWorkflowTaskSlackCanvasActionCheck(action, address, name string, updated bool) resource.TestCheckFunc {
	if action == "create_slack_canvas" {
		title := name + " {{ incident.sequential_id }}"
		if updated {
			title = "Updated " + title
		}
		return resource.TestCheckResourceAttr(address, "task_params.0.title", title)
	}
	operation := "insert_at_end"
	if updated {
		operation = "managed_sections"
	}
	return resource.TestCheckResourceAttr(address, "task_params.0.operation", operation)
}

func canvasAcceptanceAPIURL() string {
	if url := os.Getenv("ROOTLY_API_URL"); url != "" {
		return url
	}
	return "https://api.rootly.com"
}
