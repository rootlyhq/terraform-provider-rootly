package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStatusPageTemplate(t *testing.T) {
	resName := "rootly_status_page_template.test"
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	templateTitle := acctest.RandomWithPrefix("tf-template")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageTemplateConfig(
					statusPageName,
					templateTitle,
					"Initial Update Title",
					"Initial body text",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "title", templateTitle),
					resource.TestCheckResourceAttr(resName, "update_title", "Initial Update Title"),
					resource.TestCheckResourceAttr(resName, "body", "Initial body text"),
					resource.TestCheckResourceAttr(resName, "kind", "normal"),
				),
			},
			{
				Config: testAccResourceStatusPageTemplateConfig(
					statusPageName,
					templateTitle+"-updated",
					"Updated Update Title",
					"Updated body text",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "title", templateTitle+"-updated"),
					resource.TestCheckResourceAttr(resName, "update_title", "Updated Update Title"),
					resource.TestCheckResourceAttr(resName, "body", "Updated body text"),
					resource.TestCheckResourceAttr(resName, "kind", "normal"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceStatusPageTemplateWithLiquidVariables(t *testing.T) {
	resName := "rootly_status_page_template.test"
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	templateTitle := acctest.RandomWithPrefix("tf-template")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageTemplateConfig(
					statusPageName,
					templateTitle,
					"Incident Update: {{ incident.title }}",
					"We are investigating an issue with {{ incident.title }}",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "title", templateTitle),
					resource.TestCheckResourceAttr(resName, "update_title", "Incident Update: {{ incident.title }}"),
					resource.TestCheckResourceAttr(resName, "body", "We are investigating an issue with {{ incident.title }}"),
				),
			},
		},
	})
}

func testAccResourceStatusPageTemplateConfig(statusPageName, templateTitle, updateTitle, body string) string {
	return `
resource "rootly_status_page" "test" {
	title = "` + statusPageName + `"
}

resource "rootly_status_page_template" "test" {
	status_page_id = rootly_status_page.test.id
	title          = "` + templateTitle + `"
	update_title   = "` + updateTitle + `"
	body           = "` + body + `"
	kind           = "normal"
}
`
}
