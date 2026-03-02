package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStatusPage(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-status-page")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "title", rName),
					resource.TestCheckResourceAttrSet("rootly_status_page.test", "slug"),
				),
			},
			{
				Config: testAccResourceStatusPageUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "title", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "description", "A description"),
				),
			},
		},
	})
}

func TestAccResourceStatusPage_sectionOrderDefault(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-status-page")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "title", rName),
					resource.TestCheckResourceAttrSet("rootly_status_page.test", "section_order.#"),
				),
			},
			{
				Config:   testAccResourceStatusPageConfig(rName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceStatusPage_sectionOrderExplicit(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-status-page")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageSectionOrderConfig(rName, `["incidents", "maintenance", "system_status"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "title", rName),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.#", "3"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.0", "incidents"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.1", "maintenance"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.2", "system_status"),
				),
			},
			{
				Config:   testAccResourceStatusPageSectionOrderConfig(rName, `["incidents", "maintenance", "system_status"]`),
				PlanOnly: true,
			},
			{
				Config: testAccResourceStatusPageSectionOrderConfig(rName, `["system_status", "incidents", "maintenance"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.0", "system_status"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.1", "incidents"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "section_order.2", "maintenance"),
				),
			},
		},
	})
}

func testAccResourceStatusPageConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_status_page" "test" {
  title = "%s"
}
`, name)
}

func testAccResourceStatusPageUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_status_page" "test" {
  title       = "%s-updated"
  description = "A description"
}
`, name)
}

func testAccResourceStatusPageSectionOrderConfig(name, sectionOrder string) string {
	return fmt.Sprintf(`
resource "rootly_status_page" "test" {
  title         = "%s"
  section_order = %s
}
`, name, sectionOrder)
}
