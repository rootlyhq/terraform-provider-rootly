package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowPostMortem(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-pm")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowPostMortemConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "trigger_params.0.incident_visibilities.0", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowPostMortemUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "enabled", "false"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_post_mortem.foo", "trigger_params.0.incident_visibilities.0", "false"),
				),
			},
		},
	})
}

func testAccResourceWorkflowPostMortemConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_post_mortem" "foo" {
  name = "%s"
	trigger_params {
		triggers = ["post_mortem_created"]
		incident_visibilities = [true]
	}
}
`, name)
}

func testAccResourceWorkflowPostMortemUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_post_mortem" "foo" {
  name       = "%s"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["post_mortem_updated"]
		incident_visibilities = [false]
	}
}
`, name)
}
