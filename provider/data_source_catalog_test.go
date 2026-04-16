package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCatalog(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-catalog")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatalogConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_catalog.test", "name", rName),
					resource.TestCheckResourceAttrPair(
						"data.rootly_catalog.test", "id",
						"rootly_catalog.test", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceCatalogConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s"
}

data "rootly_catalog" "test" {
  name       = "%s"
  depends_on = [rootly_catalog.test]
}
`, name, name)
}
