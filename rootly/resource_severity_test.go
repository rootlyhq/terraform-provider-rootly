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
