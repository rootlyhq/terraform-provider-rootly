package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestStatusPageAnnouncementSchemaContracts(t *testing.T) {
	fields := resourceStatusPageAnnouncement().Schema
	for _, name := range []string{"user_id", "published_at"} {
		field := fields[name]
		if !field.Computed || field.Optional || field.Required {
			t.Errorf("%s must be read-only, got Computed=%t Optional=%t Required=%t", name, field.Computed, field.Optional, field.Required)
		}
	}

	notify := fields["notify_subscribers"]
	if !notify.Optional || !notify.ForceNew || notify.Default != true {
		t.Errorf("notify_subscribers must be create-only and default true, got Optional=%t ForceNew=%t Default=%v", notify.Optional, notify.ForceNew, notify.Default)
	}
}

func TestAccResourceStatusPageAnnouncement(t *testing.T) {
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	resourceName := "rootly_status_page_announcement.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageAnnouncementConfig(statusPageName, "Initial announcement", "Initial body", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Initial announcement"),
					resource.TestCheckResourceAttr(resourceName, "body", "Initial body"),
					resource.TestCheckResourceAttr(resourceName, "notify_subscribers", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "status_page_id", "rootly_status_page.test", "id"),
				),
			},
			{
				Config: testAccResourceStatusPageAnnouncementConfig(statusPageName, "Updated announcement", "Updated body", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Updated announcement"),
					resource.TestCheckResourceAttr(resourceName, "body", "Updated body"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"notify_subscribers"},
			},
		},
	})
}

func testAccResourceStatusPageAnnouncementConfig(statusPageName, title, body string, notifySubscribers bool) string {
	return fmt.Sprintf(`
resource "rootly_status_page" "test" {
  title = %q
}

resource "rootly_status_page_announcement" "test" {
  status_page_id = rootly_status_page.test.id
  title          = %q
  body           = %q
	  notify_subscribers = %t
}
`, statusPageName, title, body, notifySubscribers)
}
