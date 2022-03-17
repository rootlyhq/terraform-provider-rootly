package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCause(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcecause,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.foo", "name", "mycause"),
					resource.TestCheckResourceAttr("rootly_cause.foo", "description", ""),
				),
			},
			{
				Config: testAccResourcecauseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.foo", "name", "mycause2"),
					resource.TestCheckResourceAttr("rootly_cause.foo", "description", "my cause description"),
				),
			},
		},
	})
}

const testAccResourcecause = `
resource "rootly_cause" "foo" {
  name = "mycause"
}
`

const testAccResourcecauseUpdate = `
resource "rootly_cause" "foo" {
  name        = "mycause2"
  description = "my cause description"
}
`
