package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStatusPageAnnouncement(t *testing.T) {
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	resourceName := "rootly_status_page_announcement.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageAnnouncementConfig(statusPageName, "Initial announcement", "Initial body"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Initial announcement"),
					resource.TestCheckResourceAttr(resourceName, "body", "Initial body"),
					resource.TestCheckResourceAttrPair(resourceName, "status_page_id", "rootly_status_page.test", "id"),
				),
			},
			{
				Config: testAccResourceStatusPageAnnouncementConfig(statusPageName, "Updated announcement", "Updated body"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Updated announcement"),
					resource.TestCheckResourceAttr(resourceName, "body", "Updated body"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceStatusPageAnnouncementConfig(statusPageName, title, body string) string {
	return fmt.Sprintf(`
resource "rootly_status_page" "test" {
  title = %q
}

resource "rootly_status_page_announcement" "test" {
  status_page_id = rootly_status_page.test.id
  title          = %q
  body           = %q
}
`, statusPageName, title, body)
}
