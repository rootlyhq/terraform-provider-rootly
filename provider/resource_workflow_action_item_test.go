package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowActionItem(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowActionItem,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "name", "test-action-item-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowActionItemUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "name", "test-action-item-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowActionItem = `
resource "rootly_workflow_action_item" "foo" {
  name = "test-action-item-workflow"
	trigger_params {
		triggers = ["action_item_created"]
	}
}
`

const testAccResourceWorkflowActionItemUpdate = `
resource "rootly_workflow_action_item" "foo" {
  name       = "test-action-item-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["action_item_updated"]
	}
}
`
