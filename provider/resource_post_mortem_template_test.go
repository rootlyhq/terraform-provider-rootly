package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePostmortemTemplate(t *testing.T) {
	postMortemTemplateName := acctest.RandomWithPrefix("tf-post-mortem-template")
	postMortemTemplateContent := acctest.RandomWithPrefix("tf-post-mortem-content")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePostmortemTemplateConfig(postMortemTemplateName, postMortemTemplateContent),
			},
		},
	})
}

func testAccResourcePostmortemTemplateConfig(postMortemTemplateName, postMortemTemplateContent string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
	name = "%s"
	content = "%s"

	lifecycle {
		ignore_changes = [content]
	}
}
`, postMortemTemplateName, postMortemTemplateContent)
}
