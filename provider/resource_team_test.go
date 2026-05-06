package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTeam(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_team.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_team.test", "description", "Test team"),
				),
			},
			{
				Config: testAccResourceTeamConfigUpdated(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_team.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_team.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccResourceTeamConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name        = "%s"
	description = "Test team"
}
`, name)
}

func testAccResourceTeamConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name        = "%s"
	description = "Updated description"
}
`, name)
}
