package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	userEmail := "bot-tftests+1@rootly.com"
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserConfig(userEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_user.test", "email", userEmail),
				),
			},
		},
	})
}

func testAccDataSourceUserConfig(userEmail string) string {
	return fmt.Sprintf(`
	data "rootly_user" "test" {
		email = "%s"
	}
`, userEmail)
}
