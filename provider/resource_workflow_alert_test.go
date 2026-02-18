package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowAlert(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", "test-alert-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowAlertUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", "test-alert-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowAlert = `
resource "rootly_workflow_alert" "foo" {
  name = "test-alert-workflow"
	trigger_params {
		triggers = ["alert_created"]
	}
}
`

const testAccResourceWorkflowAlertUpdate = `
resource "rootly_workflow_alert" "foo" {
  name       = "test-alert-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["alert_created"]
	}
}
`

func TestAccResourceWorkflowAlertWithPayloadConditions(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlertWithPayloadConditions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "name", "test-alert-payload-conditionsTest"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.triggers.0", "alert_created"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.logic", "ALL"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.query", "$.commonLabels.namespace"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.operator", "IS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.values.0", "production"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.1.query", "$.commonLabels.severity"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.1.operator", "CONTAINS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.1.values.0", "critical"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.1.values.1", "high"),
				),
			},
			{
				Config: testAccResourceWorkflowAlertWithPayloadConditionsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "name", "test-alert-payload-conditions-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.logic", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.query", "$.commonLabels.environment"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.operator", "IS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.values.0", "staging"),
				),
			},
		},
	})
}

const testAccResourceWorkflowAlertWithPayloadConditions = `
resource "rootly_workflow_alert" "test_payload" {
  name    = "test-alert-payload-conditionsTest"
  enabled = true

  trigger_params {
    triggers = ["alert_created"]

    alert_payload_conditions {
      logic = "ALL"

      conditions {
        query    = "$.commonLabels.namespace"
        operator = "IS"
        values   = ["production"]
      }

      conditions {
        query    = "$.commonLabels.severity"
        operator = "CONTAINS"
        values   = ["critical", "high"]
      }
    }
  }
}
`

const testAccResourceWorkflowAlertWithPayloadConditionsUpdate = `
resource "rootly_workflow_alert" "test_payload" {
  name    = "test-alert-payload-conditions-updated"
  enabled = true

  trigger_params {
    triggers = ["alert_created"]

    alert_payload_conditions {
      logic = "ANY"

      conditions {
        query    = "$.commonLabels.environment"
        operator = "IS"
        values   = ["staging"]
      }
    }
  }
}
`

func TestAccResourceWorkflowAlertWithRegexpConditions(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlertWithRegexpConditions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "name", "test-alert-regexp-conditions"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.query", "$.alertname"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.operator", "IS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.values.0", "^(api|web)-.+"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.use_regexp", "true"),
				),
			},
		},
	})
}

const testAccResourceWorkflowAlertWithRegexpConditions = `
resource "rootly_workflow_alert" "test_regexp" {
  name    = "test-alert-regexp-conditions"
  enabled = true

  trigger_params {
    triggers = ["alert_created"]

    alert_payload_conditions {
      logic = "ALL"

      conditions {
        query      = "$.alertname"
        operator   = "IS"
        values     = ["^(api|web)-.+"]
        use_regexp = true
      }
    }
  }
}
`
