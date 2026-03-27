package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCatalogChecklistTemplate(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-checklist-tpl")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatalogChecklistTemplateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_catalog_checklist_template.test", "name", rName),
					resource.TestCheckResourceAttrPair(
						"data.rootly_catalog_checklist_template.test", "id",
						"rootly_catalog_checklist_template.test", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceCatalogChecklistTemplateConfig(name string) string {
	return fmt.Sprintf(`
data "rootly_user" "test" {
}

resource "rootly_catalog" "test" {
  name = "%s-catalog"
}

resource "rootly_catalog_checklist_template" "test" {
  catalog_type = "Catalog"
  scope_type = "Catalog"
  scope_id = rootly_catalog.test.id
  name       = "%s"
  owners {
  	id = data.rootly_user.test.id
  	type = "user"
  }
  fields {
  	field_source = "builtin"
  	field_key = "description"
  }
}

data "rootly_catalog_checklist_template" "test" {
  name       = "%s"
  depends_on = [rootly_catalog_checklist_template.test]
}
`, name, name, name)
}
