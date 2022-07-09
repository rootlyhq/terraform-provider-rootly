package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFunctionality(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
			time.Sleep(1 * time.Second)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFunctionality,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_functionality.foo", "name", "myfunctionality"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "description", ""),
					resource.TestCheckResourceAttrSet("rootly_functionality.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "color", "#047BF8"),
				),
			},
			{
				Config: testAccResourceFunctionalityUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_functionality.foo", "name", "myfunctionality2"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "description", "my functionality description"),
					resource.TestCheckResourceAttrSet("rootly_functionality.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_functionality.foo", "color", "#203"),
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
  color       = "#203"
}
`
