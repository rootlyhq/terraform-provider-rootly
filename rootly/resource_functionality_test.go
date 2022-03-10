package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFunctionality(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFunctionality,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_functionality.foo", "name", "myfunctionality"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "description", ""),
				),
			},
			{
				Config: testAccResourceFunctionalityUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_functionality.foo", "name", "myfunctionality2"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "description", "my functionality description"),
				),
			},
		},
	})
}

const testAccResourceFunctionality = `
resource "rootly_functionality" "foo" {
  name = "myfunctionality"
}
`

const testAccResourceFunctionalityUpdate = `
resource "rootly_functionality" "foo" {
  name        = "myfunctionality2"
  description = "my functionality description"
}
`
