package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEscalationPath(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPath,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "name", "test-path"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "default", "false"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restriction_time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.#", "2"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_time", "17:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_day", "tuesday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_time", "07:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.1.start_time", "17:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.1.end_day", "wednesday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.1.end_time", "07:00"),
				),
			},
			{
				Config: testAccResourceEscalationPathUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "name", "test-path-updated"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "default", "false"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restriction_time_zone", "Pacific/Honolulu"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.#", "1"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_day", "friday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_time", "18:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_time", "08:00"),
				),
			},
		},
	})
}

const testAccResourceEscalationPath = `
resource "rootly_escalation_policy" "test" {
	name = "test-ep"
}

resource "rootly_escalation_path" "test" {
	name = "test-path"
	default = false
	escalation_policy_id = rootly_escalation_policy.test.id
	time_restriction_time_zone = "America/New_York"
	time_restrictions {
		start_day = "monday"
		start_time = "17:00"
		end_day = "tuesday"
		end_time = "07:00"
	}
	time_restrictions {
		start_day = "tuesday"
		start_time = "17:00"
		end_day = "wednesday"
		end_time = "07:00"
	}
}
`

const testAccResourceEscalationPathUpdated = `
resource "rootly_escalation_policy" "test" {
	name = "test-ep"
}

resource "rootly_escalation_path" "test" {
	name = "test-path-updated"
	default = false
	escalation_policy_id = rootly_escalation_policy.test.id
	time_restriction_time_zone = "Pacific/Honolulu"
	time_restrictions {
		start_day = "friday"
		start_time = "18:00"
		end_day = "monday"
		end_time = "08:00"
	}
}
`
