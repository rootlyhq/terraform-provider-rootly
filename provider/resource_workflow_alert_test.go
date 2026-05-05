package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowAlert(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-alert")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlertConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowAlertUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "false"),
				),
			},
		},
	})
}

func testAccResourceWorkflowAlertConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_alert" "foo" {
  name = "%s"
	trigger_params {
		triggers = ["alert_created"]
	}
}
`, name)
}

func testAccResourceWorkflowAlertUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_alert" "foo" {
  name       = "%s"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["alert_created"]
	}
}
`, name)
}

func TestAccResourceWorkflowAlertWithPayloadConditions(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-alert-pc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlertPayloadConditionsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "name", rName),
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
				Config: testAccResourceWorkflowAlertPayloadConditionsUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.logic", "ANY"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.query", "$.commonLabels.environment"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.operator", "IS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_payload", "trigger_params.0.alert_payload_conditions.0.conditions.0.values.0", "staging"),
				),
			},
		},
	})
}

func testAccResourceWorkflowAlertPayloadConditionsConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_alert" "test_payload" {
  name    = "%s"
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
`, name)
}

func testAccResourceWorkflowAlertPayloadConditionsUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_alert" "test_payload" {
  name    = "%s"
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
`, name)
}

func TestAccResourceWorkflowAlertWithRegexpConditions(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-alert-re")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowAlertRegexpConditionsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.query", "$.alertname"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.operator", "IS"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.values.0", "^(api|web)-.+"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.test_regexp", "trigger_params.0.alert_payload_conditions.0.conditions.0.use_regexp", "true"),
				),
			},
		},
	})
}

func testAccResourceWorkflowAlertRegexpConditionsConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_alert" "test_regexp" {
  name    = "%s"
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
`, name)
}
