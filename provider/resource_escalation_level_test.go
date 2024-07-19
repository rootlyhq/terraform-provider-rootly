package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEscalationLevel(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationLevel,
			},
		},
	})
}

const testAccResourceEscalationLevel = `
data "rootly_user" "test" {}

resource "rootly_escalation_policy" "test" {
	name = "test"
}

resource "rootly_escalation_level" "test" {
	escalation_policy_id = rootly_escalation_policy.test.id
	position = 1
	notification_target_params {
		id = data.rootly_user.test.id
		type = "user"
	}
}
`
