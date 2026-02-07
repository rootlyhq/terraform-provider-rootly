package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEdgeConnector(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-edge-connector")
	randomNameUpdated := acctest.RandomWithPrefix("tf-test-edge-connector-updated")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEdgeConnector(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector.test", "id"),
					resource.TestCheckResourceAttr("rootly_edge_connector.test", "data.0.attributes.0.name", randomName),
					resource.TestCheckResourceAttr("rootly_edge_connector.test", "data.0.attributes.0.status", "active"),
				),
			},
			{
				Config: testAccResourceEdgeConnectorUpdate(randomNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_edge_connector.test", "id"),
					resource.TestCheckResourceAttr("rootly_edge_connector.test", "data.0.attributes.0.name", randomNameUpdated),
					resource.TestCheckResourceAttr("rootly_edge_connector.test", "data.0.attributes.0.status", "paused"),
				),
			},
		},
	})
}

func testAccResourceEdgeConnector(name string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector" "test" {
	data {
		type = "edge_connectors"
		id   = "temp-id"
		attributes {
			name          = "%s"
			description   = "Test edge connector"
			status        = "active"
			subscriptions = ["incident.created", "incident.updated"]
		}
	}
}
`, name)
}

func testAccResourceEdgeConnectorUpdate(name string) string {
	return fmt.Sprintf(`
resource "rootly_edge_connector" "test" {
	data {
		type = "edge_connectors"
		id   = "temp-id"
		attributes {
			name          = "%s"
			description   = "Updated edge connector"
			status        = "paused"
			subscriptions = ["incident.created"]
		}
	}
}
`, name)
}
