package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceHeartbeat(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHeartbeat,
			},
		},
	})
}

const testAccResourceHeartbeat = `
resource "rootly_heartbeat" "test" {
	name = "test"
alert_summary = "test"
interval = 1
interval_unit = "seconds"
notification_target_id = "test"
notification_target_type = "User"
}
`
