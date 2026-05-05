package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateOpsgenieIncident(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-opsgenie")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateOpsgenieIncidentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.task_type", "update_opsgenie_incident"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.opsgenie_incident_id", "12345"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.priority", "P1"),
				),
			},
			{
				Config: testAccResourceWorkflowTaskUpdateOpsgenieIncidentEmptyPriorityConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_opsgenie_incident.foo", "task_params.0.priority", ""),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskUpdateOpsgenieIncidentConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_opsgenie_incident" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    opsgenie_incident_id = "12345"
    priority             = "P1"
  }
}
`, name)
}

func testAccResourceWorkflowTaskUpdateOpsgenieIncidentEmptyPriorityConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_opsgenie_incident" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    opsgenie_incident_id = "12345"
    priority             = ""
  }
}
`, name)
}
