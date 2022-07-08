package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowPulse(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowPulse,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "name", "test-pulse-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowPulseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "name", "test-pulse-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowPulse = `
resource "rootly_workflow_pulse" "foo" {
  name = "test-pulse-workflow"
	trigger_params {
		triggers = ["pulse_created"]
	}
}
`

const testAccResourceWorkflowPulseUpdate = `
resource "rootly_workflow_pulse" "foo" {
  name       = "test-pulse-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["pulse_created"]
	}
}
`
