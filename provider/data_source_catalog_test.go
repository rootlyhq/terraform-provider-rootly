package provider

import (
	"fmt"
	"testing"
	"time"

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
				Config: testAccDataSourceCatalogResourceOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog.test", "name", rName),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Second) },
				Config:    testAccDataSourceCatalogWithLookup(rName),
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

func testAccDataSourceCatalogResourceOnly(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s"
}
`, name)
}

func testAccDataSourceCatalogWithLookup(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s"
}

data "rootly_catalog" "test" {
  name       = rootly_catalog.test.name
  depends_on = [rootly_catalog.test]
}
`, name)
}
