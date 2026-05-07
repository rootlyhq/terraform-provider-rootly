package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCatalogChecklistTemplate(t *testing.T) {
	catalogName := acctest.RandomWithPrefix("tf-catalog")
	checklistName := acctest.RandomWithPrefix("tf-checklist")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogChecklistTemplateConfig(catalogName, checklistName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_checklist_template.test", "name", checklistName),
				),
			},
		},
	})
}

func testAccResourceCatalogChecklistTemplateConfig(catalogName, checklistName string) string {
	return fmt.Sprintf(`
resource "rootly_catalog" "test" {
	name = "%s"
}

resource "rootly_catalog_checklist_template" "test" {
	catalog_type = rootly_catalog.test.slug
	name         = "%s"
}
`, catalogName, checklistName)
}
