package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEscalationPath(t *testing.T) {
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")
	escalationPathName := acctest.RandomWithPrefix("tf-escalation-path")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPathConfig(escalationPolicyName, escalationPathName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "name", escalationPathName),
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
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "initial_delay", "5"),
				),
			},
			{
				Config: testAccResourceEscalationPathConfigUpdated(escalationPolicyName, escalationPathName+"-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "name", escalationPathName+"-updated"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "default", "false"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restriction_time_zone", "Pacific/Honolulu"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.#", "1"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_day", "friday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.start_time", "18:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "time_restrictions.0.end_time", "08:00"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "initial_delay", "0"),
				),
			},
		},
	})
}

func testAccResourceEscalationPathConfig(policyName string, rootlyEscalationPathName string) string {

	return fmt.Sprintf(
		`
		resource "rootly_escalation_policy" "test" {
			name = "%s"
		}
		
		resource "rootly_escalation_path" "test" {
			name = "%s"
			default = false
			escalation_policy_id = rootly_escalation_policy.test.id
			initial_delay = 5
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
		`, policyName, rootlyEscalationPathName)
}

func testAccResourceEscalationPathConfigUpdated(policyName string, rootlyEscalationPathName string) string {
	return fmt.Sprintf(
		`
resource "rootly_escalation_policy" "test" {
	name = "%s"
}

resource "rootly_escalation_path" "test" {
	name = "%s"
	default = false
	escalation_policy_id = rootly_escalation_policy.test.id
	initial_delay = 0
	time_restriction_time_zone = "Pacific/Honolulu"
	time_restrictions {
		start_day = "friday"
		start_time = "18:00"
		end_day = "monday"
		end_time = "08:00"
	}
}
`, policyName, rootlyEscalationPathName)
}
