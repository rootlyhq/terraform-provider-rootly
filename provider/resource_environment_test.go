package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEnvironment(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceenvironment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_environment.foo", "name", "myenvironment"),
					resource.TestCheckResourceAttr("rootly_environment.foo", "description", ""),
				),
			},
			{
				Config: testAccResourceenvironmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_environment.foo", "name", "myenvironment2"),
					resource.TestCheckResourceAttr("rootly_environment.foo", "description", "my environment description"),
				),
			},
		},
	})
}

const testAccResourceenvironment = `
resource "rootly_environment" "foo" {
  name = "myenvironment"
}
`

const testAccResourceenvironmentUpdate = `
resource "rootly_environment" "foo" {
  name        = "myenvironment2"
  description = "my environment description"
}
`
