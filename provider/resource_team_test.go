package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTeam(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_team.foo", "name", "myteam"),
					resource.TestCheckResourceAttr("rootly_team.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_team.foo", "color", "#E65252"),
				),
			},
			{
				Config: testAccResourceTeamUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_team.foo", "name", "myteam2"),
					resource.TestCheckResourceAttr("rootly_team.foo", "description", "mydesc"),
					resource.TestCheckResourceAttr("rootly_team.foo", "color", "#123"),
				),
			},
		},
	})
}

const testAccResourceTeam = `
resource "rootly_team" "foo" {
  name = "myteam"
}
`

const testAccResourceTeamUpdate = `
resource "rootly_team" "foo" {
  name        = "myteam2"
  description = "mydesc"
  color       = "#123"
}
`
