package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFunctionality(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-functionality")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFunctionalityConfig(rName),
			},
		},
	})
}

func testAccResourceFunctionalityConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_functionality" "test" {
	name = "%s"
}
`, name)
}
