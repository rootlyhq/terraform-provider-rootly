package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWebhooksEndpoint(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-webhooks-endpoint")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWebhooksEndpointConfig(rName),
			},
		},
	})
}

func testAccResourceWebhooksEndpointConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_webhooks_endpoint" "test" {
	name = "%s"
	url = "https://rootly.com/dummy"
}
`, name)
}
