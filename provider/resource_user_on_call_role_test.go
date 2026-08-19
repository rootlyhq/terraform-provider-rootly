package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Assigns the seed user their own current on-call role so the test is
// self-contained and does not depend on knowing a specific role id.
func TestAccResourceUserOnCallRole(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserOnCallRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"rootly_user_on_call_role.test", "user_id",
						"data.rootly_user.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"rootly_user_on_call_role.test", "on_call_role_id",
						"data.rootly_user.test", "on_call_role_id",
					),
				),
			},
		},
	})
}

const testAccResourceUserOnCallRoleConfig = `
	data "rootly_user" "test" {
		email = "bot+tftests@rootly.com"
	}

	resource "rootly_user_on_call_role" "test" {
		user_id         = data.rootly_user.test.id
		on_call_role_id = data.rootly_user.test.on_call_role_id
	}
`
