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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_post_mortem_template.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_post_mortem_template.test", "format", "html"),
				),
			},
			{
				Config: testAccResourcePostmortemTemplateConfigUpdated(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_post_mortem_template.test", "name", rName+"-updated"),
				),
			},
		},
	})
}

func TestAccResourcePostmortemTemplate_MarkdownNoDrift(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-pm-md")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePostmortemTemplateMarkdownConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_post_mortem_template.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_post_mortem_template.test", "format", "markdown"),
				),
			},
			// Second plan with same config should show no changes (no drift)
			{
				Config:   testAccResourcePostmortemTemplateMarkdownConfig(rName),
				PlanOnly: true,
			},
		},
	})
}

func testAccResourcePostmortemTemplateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
	name    = "%s"
	content = "<p>Test HTML content</p>"
	format  = "html"
}
`, name)
}

func testAccResourcePostmortemTemplateConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
	name    = "%s"
	content = "<p>Updated HTML content</p>"
	format  = "html"
}
`, name)
}

func testAccResourcePostmortemTemplateMarkdownConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_post_mortem_template" "test" {
	name    = "%s"
	content = "**Summary**\n\n*Brief summary of what happened.*\n\n{{ incident.summary }}"
	format  = "markdown"
}
`, name)
}
