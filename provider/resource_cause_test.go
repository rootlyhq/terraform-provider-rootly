package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCause(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCause,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("causes", "bug"),
				),
			},
			{
				Config: testAccResourceCause,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.foo", "name", "mycause"),
					resource.TestCheckResourceAttr("rootly_cause.foo", "description", ""),
				),
			},
			{
				Config: testAccResourceCauseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.foo", "name", "mycause2"),
					resource.TestCheckResourceAttr("rootly_cause.foo", "description", "my cause description"),
				),
			},
		},
	})
}

const testAccDataSourceCause = `
data "rootly_causes" "test" {
	slug = "bug"
}

output "causes" {
	value = data.rootly_causes.test.causes[0].slug
}
`

const testAccResourceCause = `
resource "rootly_cause" "foo" {
  name = "mycause"
}
`

const testAccResourceCauseUpdate = `
resource "rootly_cause" "foo" {
  name        = "mycause2"
  description = "my cause description"
}
`
