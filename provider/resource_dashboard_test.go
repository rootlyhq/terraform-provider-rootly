package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDashboard(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDashboard,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "name", "mydashboard"),
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "description", ""),
				),
			},
			{
				Config: testAccResourceDashboardUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "name", "mydashboard2"),
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "description", "my dashboard description"),
				),
			},
		},
	})
}

const testAccResourceDashboard = `
resource "rootly_dashboard" "foo" {
  name = "mydashboard"
}
`

const testAccResourceDashboardUpdate = `
resource "rootly_dashboard" "foo" {
  name        = "mydashboard2"
}
`
