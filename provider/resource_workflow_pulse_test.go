package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowPulse(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-pulse")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowPulseConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowPulseUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_pulse.foo", "enabled", "false"),
				),
			},
		},
	})
}

func testAccResourceWorkflowPulseConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_pulse" "foo" {
  name = "%s"
	trigger_params {
		triggers = ["pulse_created"]
	}
}
`, name)
}

func testAccResourceWorkflowPulseUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_pulse" "foo" {
  name       = "%s"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["pulse_created"]
	}
}
`, name)
}
