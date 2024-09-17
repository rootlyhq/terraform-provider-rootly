package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertUrgency(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertUrgency,
			},
		},
	})
}

const testAccResourceAlertUrgency = `
resource "rootly_alert_urgency" "test" {
	name = "test"
description = "test"
}
`
