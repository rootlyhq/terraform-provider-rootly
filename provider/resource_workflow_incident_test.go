package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.0", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowIncidentUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "name", rName+"-3"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.0", "false"),
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
		incident_visibilities = [true]
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
		incident_visibilities = [false]
	}
}
`, rName)
}
