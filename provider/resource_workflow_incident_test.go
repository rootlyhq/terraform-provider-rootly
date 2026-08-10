package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceWorkflowIncident(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-inc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowIncidentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "name", rName+"-3"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo1", "trigger_params.0.incident_visibilities.0", "true"),
				),
			},
			{
				// Regression coverage: incident_visibilities is a boolean array in the
				// API, so importing a workflow that sets it used to drop the whole
				// trigger_params block.
				ResourceName: "rootly_workflow_incident.foo1",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					if got := states[0].Attributes["trigger_params.0.incident_visibilities.0"]; got != "true" {
						return fmt.Errorf("expected imported incident_visibilities.0 to be \"true\", got %q", got)
					}
					return nil
				},
			},
			{
				Config: testAccResourceWorkflowIncidentUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "name", rName+"-3"),
				),
			},
		},
	})
}

func testAccResourceWorkflowIncidentConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo1" {
  name = "%s-1"
	trigger_params {
		triggers = ["incident_updated"]
		incident_condition_visibility = "IS"
		incident_visibilities = [true]
	}
}
resource "rootly_workflow_incident" "foo2" {
  name = "%s-2"
	trigger_params {
		triggers = ["incident_updated"]
	}
	depends_on = [rootly_workflow_incident.foo1]
}
resource "rootly_workflow_incident" "foo3" {
  name = "%s-3"
	trigger_params {
		triggers = ["incident_updated"]
	}
	depends_on =[rootly_workflow_incident.foo2]
}
`, rName, rName, rName)
}

func testAccResourceWorkflowIncidentUpdateConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo3" {
  name = "%s-3"
	trigger_params {
		triggers = ["incident_updated"]
	}
}
`, rName)
}
