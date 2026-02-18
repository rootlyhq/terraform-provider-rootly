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
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "initial_delay", "5"),
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
					resource.TestCheckResourceAttr("rootly_escalation_path.test", "initial_delay", "0"),
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
`

const testAccResourceEscalationPathUpdated = `
resource "rootly_escalation_policy" "test" {
	name = "test-ep"
}

resource "rootly_escalation_path" "test" {
	name = "test-path-updated"
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
`

func TestAccResourceEscalationPathWithAllRuleTypes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPathWithAllRuleTypes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "name", "tf-test-path-all-rules"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "match_mode", "match-any-rule"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "rules.#", "4"),
				),
			},
			{
				Config: testAccResourceEscalationPathWithAllRuleTypesUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "name", "tf-test-path-all-rules-updated"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "match_mode", "match-all-rules"),
					resource.TestCheckResourceAttr("rootly_escalation_path.test_rules", "rules.#", "2"),
				),
			},
		},
	})
}

const testAccResourceEscalationPathWithAllRuleTypes = `
resource "rootly_escalation_policy" "test_rules" {
	name = "tf-test-ep-rules"
}

resource "rootly_alert_urgency" "test" {
	name = "tf-test-urgency"
	description = "Test urgency for escalation path rules"
}

resource "rootly_alert_field" "test" {
	name = "tf-test-alert-field"
}

resource "rootly_escalation_path" "test_rules" {
	name = "tf-test-path-all-rules"
	default = false
	escalation_policy_id = rootly_escalation_policy.test_rules.id
	match_mode = "match-any-rule"

	# Rule type: alert_urgency
	rules {
		rule_type = "alert_urgency"
		urgency_ids = [rootly_alert_urgency.test.id]
	}

	# Rule type: working_hour
	rules {
		rule_type = "working_hour"
		within_working_hour = true
	}

	# Rule type: json_path
	rules {
		rule_type = "json_path"
		json_path = "$.severity"
		operator = "is"
		value = "critical"
	}

	# Rule type: field
	rules {
		rule_type = "field"
		fieldable_type = "AlertField"
		fieldable_id = rootly_alert_field.test.id
		operator = "is_one_of"
		values = ["value1", "value2"]
	}
}
`

const testAccResourceEscalationPathWithAllRuleTypesUpdated = `
resource "rootly_escalation_policy" "test_rules" {
	name = "tf-test-ep-rules"
}

resource "rootly_alert_urgency" "test" {
	name = "tf-test-urgency"
	description = "Test urgency for escalation path rules"
}

resource "rootly_alert_field" "test" {
	name = "tf-test-alert-field"
}

resource "rootly_escalation_path" "test_rules" {
	name = "tf-test-path-all-rules-updated"
	default = false
	escalation_policy_id = rootly_escalation_policy.test_rules.id
	match_mode = "match-all-rules"

	# Rule type: field with different operator
	rules {
		rule_type = "field"
		fieldable_type = "AlertField"
		fieldable_id = rootly_alert_field.test.id
		operator = "contains"
		values = ["updated-value"]
	}

	# Rule type: json_path with different operator
	rules {
		rule_type = "json_path"
		json_path = "$.priority"
		operator = "is_not"
		value = "low"
	}
}
`
