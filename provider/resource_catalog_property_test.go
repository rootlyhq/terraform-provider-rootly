package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCatalogProperty(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-catalog-prop")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogPropertyConfig(rName, "text", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "kind", "text"),
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "required", "false"),
				),
			},
			{
				Config: testAccResourceCatalogPropertyConfig(rName+"-updated", "text", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "required", "true"),
				),
			},
		},
	})
}

func TestAccResourceCatalogPropertyWithCatalogType(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-catalog-prop")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogPropertyWithCatalogTypeConfig(rName, "service"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_catalog_property.test", "catalog_type", "service"),
				),
			},
		},
	})
}

func testAccResourceCatalogPropertyConfig(name, kind string, required bool) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_property" "test" {
  catalog_id = rootly_catalog.test.id
  name       = "%s"
  kind       = "%s"
  required   = %t
}
`, name, name, kind, required)
}

func testAccResourceCatalogPropertyWithCatalogTypeConfig(name, catalogType string) string {
	return fmt.Sprintf(`
resource "rootly_catalog_property" "test" {
  name         = "%s"
  catalog_type = "%s"
}
`, name, catalogType)
}
