package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_user.test", "email", "bot+tftests@rootly.com"),
					resource.TestCheckResourceAttrSet("data.rootly_user.test", "on_call_role_id"),
				),
			},
		},
	})
}

const testAccDataSourceUser = `
	data "rootly_user" "test" {
		email = "bot+tftests@rootly.com"
	}
`
