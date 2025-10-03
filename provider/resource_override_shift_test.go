package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOverrideShift(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOverrideShift,
			},
		},
	})
}

const testAccResourceOverrideShift = `
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "primary" {
	name              = "TF Test Schedule"
	owner_user_id     = data.rootly_user.test.id
	all_time_coverage = false
}

resource "rootly_override_shift" "test" {
	schedule_id = rootly_schedule.primary.id
	starts_at   = "2026-10-01T00:00:00.000-07:00"
	ends_at     = "2026-10-02T00:00:00.000-07:00"
	user_id     = data.rootly_user.test.id
}
`
