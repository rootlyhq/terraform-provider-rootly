package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCatalogChecklistTemplate(t *testing.T) {
	checklistName := acctest.RandomWithPrefix("tf-checklist")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCatalogChecklistTemplateConfig(checklistName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_catalog_checklist_template.test", "name", checklistName),
					resource.TestCheckResourceAttr("rootly_catalog_checklist_template.test", "catalog_type", "Service"),
				),
			},
		},
	})
}

func testAccResourceCatalogChecklistTemplateConfig(checklistName string) string {
	return fmt.Sprintf(`
resource "rootly_catalog_checklist_template" "test" {
	name         = "%s"
	catalog_type = "Service"
}
`, checklistName)
}
