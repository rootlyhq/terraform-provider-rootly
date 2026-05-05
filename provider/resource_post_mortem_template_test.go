package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePostmortemTemplate(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-pm-tpl")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePostmortemTemplateConfig(rName),
			},
		},
	})
}

func testAccResourcePostmortemTemplateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
	name = "%s"
	content = "test"

	lifecycle {
		ignore_changes = [content]
	}
}
`, name)
}
