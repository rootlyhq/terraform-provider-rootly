package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRole(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-role")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_role.test", "name", rName),
				),
			},
			{
				Config: testAccResourceRoleConfigUpdated(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_role.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_role.test", "incidents_permissions.#", "2"),
				),
			},
		},
	})
}

func testAccResourceRoleConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_role" "test" {
	name = "%s"
}
`, name)
}

func testAccResourceRoleConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "rootly_role" "test" {
	name = "%s"
	incidents_permissions = ["read", "update"]
}
`, name)
}
