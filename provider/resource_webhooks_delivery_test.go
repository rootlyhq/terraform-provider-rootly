package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWebhooksDelivery(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWebhooksDelivery,
			},
		},
	})
}

const testAccResourceWebhooksDelivery = `
resource "rootly_webhooks_delivery" "test" {
	
}
`
