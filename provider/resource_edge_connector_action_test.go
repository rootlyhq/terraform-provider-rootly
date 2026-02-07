package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEdgeConnectorAction(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-edge-action")
	randomNameUpdated := acctest.RandomWithPrefix("tf-test-edge-action-updated")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEdgeConnectorAction(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector_action.test", "id"),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.name", randomName),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.action_type", "script"),
				),
			},
			{
				Config: testAccResourceEdgeConnectorActionUpdate(randomNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector_action.test", "id"),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.name", randomNameUpdated),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.action_type", "http"),
				),
			},
		},
	})
}

func testAccResourceEdgeConnectorAction(name string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector_action" "test" {
	data {
		type = "edge_connector_actions"
		id   = "temp-id"
		attributes {
			name        = "%s"
			description = "Test edge connector action"
			action_type = "script"
			icon        = "bolt"
			timeout     = 300
			parameters {
				name        = "test_param"
				type        = "string"
				required    = true
				description = "Test parameter"
			}
		}
	}
}
`, name)
}

func testAccResourceEdgeConnectorActionUpdate(name string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector_action" "test" {
	data {
		type = "edge_connector_actions"
		id   = "temp-id"
		attributes {
			name        = "%s"
			description = "Updated edge connector action"
			action_type = "http"
			icon        = "rocket-launch"
			timeout     = 600
			parameters {
				name        = "updated_param"
				type        = "number"
				required    = false
				description = "Updated parameter"
			}
		}
	}
}
`, name)
}
