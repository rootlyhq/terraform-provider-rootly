package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAddRole(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-add-role")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAddRoleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_add_role.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_role.foo", "task_params.0.task_type", "add_role"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskAddRoleConfig(name string) string {
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

resource "rootly_workflow_task_add_role" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    incident_role_id = rootly_incident_role.foo.id
  }
}
`, name, name)
}
