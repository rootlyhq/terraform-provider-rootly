package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// A known-good email plus one that does not exist: the missing
				// one must be omitted rather than causing an error.
				Config: testAccDataSourceUsersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.rootly_users.test", "users.0.email", "bot+tftests@rootly.com"),
				),
			},
		},
	})
}

const testAccDataSourceUsersConfig = `
	data "rootly_users" "test" {
		emails = ["bot+tftests@rootly.com", "definitely-not-a-real-user@rootly.com"]
	}
`
