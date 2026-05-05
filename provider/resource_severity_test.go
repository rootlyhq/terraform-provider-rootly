package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSeverity(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-severity")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSeverityConfig(rName),
			},
		},
	})
}

func testAccResourceSeverityConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_severity" "test" {
	name = "%s"
}
`, name)
}
