package provider

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr("rootly_team.test", "schedule_override_policy", "everyone"),
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

func TestAccResourceTeam_ScheduleOverridePolicy(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-team-sop")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfigScheduleOverridePolicy(rName, "admins"),
				Check:  resource.TestCheckResourceAttr("rootly_team.test", "schedule_override_policy", "admins"),
			},
			{
				Config: testAccResourceTeamConfigScheduleOverridePolicy(rName, "members"),
				Check:  resource.TestCheckResourceAttr("rootly_team.test", "schedule_override_policy", "members"),
			},
			{
				ResourceName:      "rootly_team.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTeam_AdminIdsValidation(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-team-adm")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceTeamAdminNotInUsersConfig(rName),
				ExpectError: regexp.MustCompile(`admin_ids contains user .* which is not in user_ids`),
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

func testAccResourceTeamConfigScheduleOverridePolicy(name, policy string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name                     = "%s"
	description              = "Test team"
	schedule_override_policy = "%s"
}
`, name, policy)
}

func testAccResourceTeamAdminNotInUsersConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name      = "%s"
	user_ids  = [12345]
	admin_ids = [12345, 99999]
}
`, name)
}
