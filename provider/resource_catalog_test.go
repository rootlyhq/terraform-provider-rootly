package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCatalog(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-catalog")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogConfig(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog.test", "name", rName),
				),
			},
			{
				Config: testAccResourceCatalogConfig(rName+"-updated", "A description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_catalog.test", "description", "A description"),
				),
			},
		},
	})
}

func testAccResourceCatalogConfig(name, description string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}
