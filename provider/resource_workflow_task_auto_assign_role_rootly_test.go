package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAutoAssignRoleRootly(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-aar-rt")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAutoAssignRoleRootlyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_auto_assign_role_rootly.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_auto_assign_role_rootly.foo", "task_params.0.task_type", "auto_assign_role_rootly"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskAutoAssignRoleRootlyConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "foo" {
  name = "%s-role"
}

resource "rootly_team" "foo" {
  name = "%s-team"
}

resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_auto_assign_role_rootly" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    incident_role_id = rootly_incident_role.foo.id
    group_target = {
      id   = rootly_team.foo.id
      name = rootly_team.foo.name
    }
  }
}
`, name, name, name)
}
