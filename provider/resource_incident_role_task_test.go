package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentRoleTask(t *testing.T) {
	roleName := acctest.RandomWithPrefix("tf-role")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentRoleTaskConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_role_task.test", "task", "Open a war room"),
				),
			},
		},
	})
}

func testAccResourceIncidentRoleTaskConfig(roleName string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "test" {
	name = "%s"
}

resource "rootly_incident_role_task" "test" {
	incident_role_id = rootly_incident_role.test.id
	task             = "Open a war room"
}
`, roleName)
}
