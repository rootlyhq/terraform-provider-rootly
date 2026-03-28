package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCatalogEntity(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-catalog-entity")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogEntityConfig(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "name", rName),
				),
			},
			{
				Config: testAccResourceCatalogEntityConfig(rName+"-updated", "A description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "description", "A description"),
				),
			},
		},
	})
}

func TestAccResourceCatalogEntityWithProperties(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-catalog-entity")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogEntityWithPropertiesConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_catalog_entity.test", "properties.0.value", "test-value"),
				),
			},
		},
	})
}

func testAccResourceCatalogEntityConfig(name, description string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_entity" "test" {
  catalog_id  = rootly_catalog.test.id
  name        = "%s"
  description = "%s"
}
`, name, name, description)
}

func testAccResourceCatalogEntityWithPropertiesConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_property" "test" {
  catalog_id = rootly_catalog.test.id
  name       = "%s-property"
  kind       = "text"
}

resource "rootly_catalog_entity" "test" {
  catalog_id = rootly_catalog.test.id
  name       = "%s"
  properties {
    catalog_property_id = rootly_catalog_property.test.id
    value               = "test-value"
  }
}
`, name, name, name)
}
