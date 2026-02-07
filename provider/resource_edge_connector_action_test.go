package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEdgeConnectorAction(t *testing.T) {
	randomConnectorName := acctest.RandomWithPrefix("tf-test-edge-connector")
	randomActionName := acctest.RandomWithPrefix("tf-test-edge-action")
	randomActionNameUpdated := acctest.RandomWithPrefix("tf-test-edge-action-updated")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEdgeConnectorAction(randomConnectorName, randomActionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector_action.test", "id"),
					resource.TestCheckResourceAttrSet("rootly_edge_connector_action.test", "edge_connector_id"),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.name", randomActionName),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.action_type", "script"),
				),
			},
			{
				Config: testAccResourceEdgeConnectorActionUpdate(randomConnectorName, randomActionNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector_action.test", "id"),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.name", randomActionNameUpdated),
					resource.TestCheckResourceAttr("rootly_edge_connector_action.test", "data.0.attributes.0.action_type", "http"),
				),
			},
		},
	})
}

func testAccResourceEdgeConnectorAction(connectorName string, actionName string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector" "parent" {
	data {
		type = "edge_connectors"
		id   = "temp-id"
		attributes {
			name        = "%s"
			description = "Parent edge connector for testing"
			status      = "active"
		}
	}
}

resource "rootly_edge_connector_action" "test" {
	edge_connector_id = rootly_edge_connector.parent.id
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
`, connectorName, actionName)
}

func testAccResourceEdgeConnectorActionUpdate(connectorName string, actionName string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector" "parent" {
	data {
		type = "edge_connectors"
		id   = "temp-id"
		attributes {
			name        = "%s"
			description = "Parent edge connector for testing"
			status      = "active"
		}
	}
}

resource "rootly_edge_connector_action" "test" {
	edge_connector_id = rootly_edge_connector.parent.id
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
`, connectorName, actionName)
}
