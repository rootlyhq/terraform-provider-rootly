package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStatusPageComponent(t *testing.T) {
	resName := "rootly_status_page_component.test"
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	groupName := acctest.RandomWithPrefix("tf-component-group")
	componentName := acctest.RandomWithPrefix("tf-component")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageComponentConfig(
					statusPageName,
					groupName,
					componentName,
					"Initial description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", componentName),
					resource.TestCheckResourceAttr(resName, "description", "Initial description"),
					resource.TestCheckResourceAttrSet(resName, "position"),
					resource.TestCheckResourceAttrPair(resName, "status_page_id", "rootly_status_page.test", "id"),
					resource.TestCheckResourceAttrPair(resName, "status_page_component_group_id", "rootly_status_page_component_group.test", "id"),
				),
			},
			{
				Config: testAccResourceStatusPageComponentConfig(
					statusPageName,
					groupName,
					componentName+"-updated",
					"Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", componentName+"-updated"),
					resource.TestCheckResourceAttr(resName, "description", "Updated description"),
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

func TestAccResourceStatusPageComponentCatalogBacked(t *testing.T) {
	resName := "rootly_status_page_component.test"
	statusPageName := acctest.RandomWithPrefix("tf-status-page")
	serviceName := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageComponentCatalogBackedConfig(
					statusPageName,
					serviceName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "source_type", "Service"),
					resource.TestCheckResourceAttrPair(resName, "source_id", "rootly_service.test", "id"),
					// Name is derived from the backing service for catalog-backed components.
					resource.TestCheckResourceAttr(resName, "name", serviceName),
					resource.TestCheckResourceAttrSet(resName, "position"),
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

func testAccResourceStatusPageComponentConfig(statusPageName, groupName, componentName, description string) string {
	return `
resource "rootly_status_page" "test" {
	title = "` + statusPageName + `"
}

resource "rootly_status_page_component_group" "test" {
	status_page_id = rootly_status_page.test.id
	name           = "` + groupName + `"
}

resource "rootly_status_page_component" "test" {
	status_page_id                 = rootly_status_page.test.id
	status_page_component_group_id = rootly_status_page_component_group.test.id
	name                           = "` + componentName + `"
	description                    = "` + description + `"
}
`
}

func testAccResourceStatusPageComponentCatalogBackedConfig(statusPageName, serviceName string) string {
	return `
resource "rootly_status_page" "test" {
	title = "` + statusPageName + `"
}

resource "rootly_service" "test" {
	name = "` + serviceName + `"
}

resource "rootly_status_page_component" "test" {
	status_page_id = rootly_status_page.test.id
	source_type    = "Service"
	source_id      = rootly_service.test.id
}
`
}
