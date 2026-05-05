package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowActionItem(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-action")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowActionItemConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowActionItemUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_action_item.foo", "enabled", "false"),
				),
			},
		},
	})
}

func testAccResourceWorkflowActionItemConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_action_item" "foo" {
  name = "%s"
	trigger_params {
		triggers = ["action_item_created"]
	}
}
`, name)
}

func testAccResourceWorkflowActionItemUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_action_item" "foo" {
  name       = "%s"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["action_item_updated"]
	}
}
`, name)
}
