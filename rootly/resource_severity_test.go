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
					resource.TestCheckResourceAttr("rootly_severity.foo", "name", "mysev"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "description", ""),
					resource.TestCheckResourceAttrSet("rootly_severity.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "color", "#047BF8"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "severity", "medium"),
				),
			},
			{
				Config: testAccResourceSeverityUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_severity.foo", "name", "mysev2"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "description", "test description"),
					resource.TestCheckResourceAttrSet("rootly_severity.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "color", "#203"),
					resource.TestCheckResourceAttr("rootly_severity.foo", "severity", "high"),
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
  color       = "#203"
}
`
