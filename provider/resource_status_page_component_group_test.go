package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStatusPageComponentGroup(t *testing.T) {
	resName := "rootly_status_page_component_group.test"
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	groupName := acctest.RandomWithPrefix("tf-component-group")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageComponentGroupConfig(
					statusPageName,
					groupName,
					"Initial description",
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", groupName),
					resource.TestCheckResourceAttr(resName, "description", "Initial description"),
					resource.TestCheckResourceAttr(resName, "collapsed_by_default", "true"),
					resource.TestCheckResourceAttrSet(resName, "position"),
					resource.TestCheckResourceAttrPair(resName, "status_page_id", "rootly_status_page.test", "id"),
				),
			},
			{
				Config: testAccResourceStatusPageComponentGroupConfig(
					statusPageName,
					groupName+"-updated",
					"Updated description",
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", groupName+"-updated"),
					resource.TestCheckResourceAttr(resName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resName, "collapsed_by_default", "false"),
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

func testAccResourceStatusPageComponentGroupConfig(statusPageName, groupName, description string, collapsedByDefault bool) string {
	collapsed := "false"
	if collapsedByDefault {
		collapsed = "true"
	}
	return `
resource "rootly_status_page" "test" {
	title = "` + statusPageName + `"
}

resource "rootly_status_page_component_group" "test" {
	status_page_id       = rootly_status_page.test.id
	name                 = "` + groupName + `"
	description          = "` + description + `"
	collapsed_by_default = ` + collapsed + `
}
`
}
