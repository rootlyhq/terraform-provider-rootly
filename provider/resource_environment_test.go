package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEnvironment(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-environment")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEnvironmentConfig(rName),
			},
		},
	})
}

func testAccResourceEnvironmentConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_environment" "test" {
	name = "%s"
}
`, name)
}
