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
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", "test-postmortem-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowPostMortemUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", "test-postmortem-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowPostMortem = `
resource "rootly_workflow_post_mortem" "foo" {
  name = "test-postmortem-workflow"
	trigger_params {
		triggers = ["post_mortem_created"]
	}
}
`

const testAccResourceWorkflowPostMortemUpdate = `
resource "rootly_workflow_post_mortem" "foo" {
  name       = "test-postmortem-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["post_mortem_updated"]
	}
}
`

func TestAccResourceWorkflowPostMortemWithIsNotCondition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowPostMortemWithIsNotCondition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.test_is_not", "name", "test-postmortem-workflow-is-not"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.test_is_not", "trigger_params.0.incident_condition_status", "IS NOT"),
				),
			},
		},
	})
}

const testAccResourceWorkflowPostMortemWithIsNotCondition = `
resource "rootly_workflow_post_mortem" "test_is_not" {
  name = "test-postmortem-workflow-is-not"
	trigger_params {
		triggers = ["post_mortem_created"]
		incident_condition_status = "IS NOT"
	}
}
`
