package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCatalogEntity(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-catalog-entity")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatalogEntityResourceOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "name", rName),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Second) },
				Config:    testAccDataSourceCatalogEntityWithLookup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_catalog_entity.test", "name", rName),
					resource.TestCheckResourceAttrPair(
						"data.rootly_catalog_entity.test", "id",
						"rootly_catalog_entity.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.rootly_catalog_entity.test", "catalog_id",
						"rootly_catalog.test", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceCatalogEntityResourceOnly(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_entity" "test" {
  catalog_id = rootly_catalog.test.id
  name       = "%s"
}
`, name, name)
}

func testAccDataSourceCatalogEntityWithLookup(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_entity" "test" {
  catalog_id = rootly_catalog.test.id
  name       = "%s"
}

data "rootly_catalog_entity" "test" {
  catalog_id = rootly_catalog.test.id
  name       = rootly_catalog_entity.test.name
  depends_on = [rootly_catalog_entity.test]
}
`, name, name)
}
