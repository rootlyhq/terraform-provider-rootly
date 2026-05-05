package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDashboard(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-dashboard")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "name", rName),
				),
			},
			{
				Config: testAccResourceDashboardConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_dashboard.foo", "name", rName+"-updated"),
				),
			},
		},
	})
}

func testAccResourceDashboardConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_dashboard" "foo" {
  name = "%s"
}
`, name)
}
