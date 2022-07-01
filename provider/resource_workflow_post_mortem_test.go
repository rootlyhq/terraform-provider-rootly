package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowPostMortem(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowPostMortem,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", "test-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowPostMortemUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", "test-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowPostMortem = `
resource "rootly_workflow_post_mortem" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["post_mortem_created"]
	}
}
`

const testAccResourceWorkflowPostMortemUpdate = `
resource "rootly_workflow_post_mortem" "foo" {
  name       = "test-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["post_mortem_updated"]
	}
}
`
