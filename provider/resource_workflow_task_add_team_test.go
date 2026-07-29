package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAddTeam(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-add-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAddTeamConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_add_team.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_team.foo", "task_params.0.task_type", "add_team"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskAddTeamConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_team" "foo" {
  name = "%s-team"
}

resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_add_team" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    group_id = rootly_team.foo.id
  }
}
`, name, name)
}
