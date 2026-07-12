package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePostmortemTemplate(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-pm-tpl")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePostmortemTemplateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.rootly_post_mortem_template.by_slug", "id",
						"rootly_post_mortem_template.test", "id",
					),
					resource.TestCheckResourceAttr(
						"data.rootly_post_mortem_template.by_slug", "name", name,
					),
				),
			},
		},
	})
}

func testAccDataSourcePostmortemTemplateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
  name    = "%s"
  content = "test content"

  lifecycle {
    ignore_changes = [content]
  }
}

data "rootly_post_mortem_template" "by_slug" {
  slug       = "%s"
  depends_on = [rootly_post_mortem_template.test]
}
`, name, name)
}
