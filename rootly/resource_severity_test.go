package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSeverity(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSeverity,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rootly_severity.foo", "name", "mysev"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "severity", "medium"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "description", ""),
				),
			},
			{
				Config: testAccResourceSeverityUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_severity.foo", "name", "mysev2"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "severity", "high"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "description", "test description"),
				),
			},
		},
	})
}

const testAccResourceSeverity = `
resource "rootly_severity" "foo" {
  name = "mysev"
}
`

const testAccResourceSeverityUpdate = `
resource "rootly_severity" "foo" {
  name        = "mysev2"
  severity    = "high"
  description = "test description"
}
`
