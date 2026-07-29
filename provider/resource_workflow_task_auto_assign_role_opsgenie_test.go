package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAutoAssignRoleOpsgenie(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-aar-og")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAutoAssignRoleOpsgenieConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_auto_assign_role_opsgenie.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_auto_assign_role_opsgenie.foo", "task_params.0.task_type", "auto_assign_role_opsgenie"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskAutoAssignRoleOpsgenieConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "foo" {
  name = "%s-role"
}

resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_auto_assign_role_opsgenie" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    incident_role_id = rootly_incident_role.foo.id
    schedule = {
      id   = "foo"
      name = "bar"
    }
  }
}
`, name, name)
}
