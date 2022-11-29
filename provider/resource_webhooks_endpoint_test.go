package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWebhooksEndpoint(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWebhooksEndpoint,
			},
		},
	})
}

const testAccResourceWebhooksEndpoint = `
resource "rootly_webhooks_endpoint" "test" {
	name = "test"
	url = https://rootly.com/dummy

}
`
