package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowIncident(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-inc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowIncidentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "name", rName+"-3"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.0", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowIncidentUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "name", rName+"-3"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.#", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo3", "trigger_params.0.incident_visibilities.0", "false"),
				),
			},
		},
	})
}

func testAccResourceWorkflowIncidentConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo1" {
  name = "%s-1"
	trigger_params {
		triggers = ["incident_updated"]
	}
}
resource "rootly_workflow_incident" "foo2" {
  name = "%s-2"
	trigger_params {
		triggers = ["incident_updated"]
	}
	depends_on = [rootly_workflow_incident.foo1]
}
resource "rootly_workflow_incident" "foo3" {
  name = "%s-3"
	trigger_params {
		triggers = ["incident_updated"]
		incident_visibilities = [true]
	}
	depends_on =[rootly_workflow_incident.foo2]
}
`, rName, rName, rName)
}

func testAccResourceWorkflowIncidentUpdateConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo3" {
  name = "%s-3"
	trigger_params {
		triggers = ["incident_updated"]
		incident_visibilities = [false]
	}
}
`, rName)
}

func TestAccResourceWorkflowIncidentFailureNotifications(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-inc-fn")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowIncidentFailureNotificationsCustomConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "name", rName),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_mode", "custom"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.#", "2"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.0.id", "{{ incident.slack_channel_id }}"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.0.name", "{{ incident.slack_channel_id }}"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.1.id", "C0123456789"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.1.name", "workflow-failures"),
				),
			},
			{
				Config: testAccResourceWorkflowIncidentFailureNotificationsOffConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_mode", "off"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.fn", "failure_notification_channels.#", "2"),
				),
			},
			{
				ResourceName:      "rootly_workflow_incident.fn",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceWorkflowIncidentFailureNotificationsCustomConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "fn" {
  name = "%s"
  trigger_params {
    triggers = ["incident_updated"]
  }

  failure_notification_mode = "custom"
  failure_notification_channels {
    id   = "{{ incident.slack_channel_id }}"
    name = "{{ incident.slack_channel_id }}"
  }
  failure_notification_channels {
    id   = "C0123456789"
    name = "workflow-failures"
  }
}
`, rName)
}

func testAccResourceWorkflowIncidentFailureNotificationsOffConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "fn" {
  name = "%s"
  trigger_params {
    triggers = ["incident_updated"]
  }

  failure_notification_mode = "off"
  failure_notification_channels {
    id   = "{{ incident.slack_channel_id }}"
    name = "{{ incident.slack_channel_id }}"
  }
  failure_notification_channels {
    id   = "C0123456789"
    name = "workflow-failures"
  }
}
`, rName)
}

// The deprecated incident_post_mortem_condition_cause aliases incident_condition_cause in the API,
// so it must only be sent when configured, never with its default or a stale state value.
func TestAccResourceWorkflowIncidentConditionCause(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-inc-cause")
	check := func(cause string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("rootly_workflow_incident.cause", "trigger_params.0.incident_condition_cause", cause),
			resource.TestCheckResourceAttr("rootly_workflow_incident.cause", "trigger_params.0.incident_post_mortem_condition_cause", cause),
		)
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Create with only the current attribute.
			{
				Config: testAccResourceWorkflowIncidentConditionCauseConfig(rName, true, "UNSET", ""),
				Check:  check("UNSET"),
			},
			// Update another trigger param; the condition must survive.
			{
				Config: testAccResourceWorkflowIncidentConditionCauseConfig(rName, false, "UNSET", ""),
				Check:  check("UNSET"),
			},
			// Change the condition; the deprecated attribute's state value must not be sent back.
			{
				Config: testAccResourceWorkflowIncidentConditionCauseConfig(rName, false, "SET", ""),
				Check:  check("SET"),
			},
			// Both attributes set to the same value.
			{
				Config: testAccResourceWorkflowIncidentConditionCauseConfig(rName, false, "UNSET", "UNSET"),
				Check:  check("UNSET"),
			},
			// Drop the deprecated attribute from config; the condition must stay as configured.
			{
				Config: testAccResourceWorkflowIncidentConditionCauseConfig(rName, false, "UNSET", ""),
				Check:  check("UNSET"),
			},
			{
				ResourceName:      "rootly_workflow_incident.cause",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceWorkflowIncidentConditionCauseConfig(rName string, visibility bool, cause, deprecatedCause string) string {
	deprecated := ""
	if deprecatedCause != "" {
		deprecated = fmt.Sprintf("incident_post_mortem_condition_cause = %q", deprecatedCause)
	}
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "cause" {
  name = "%s"
  trigger_params {
    triggers                 = ["incident_updated"]
    incident_visibilities    = [%t]
    incident_condition_cause = %q
    %s
  }
}
`, rName, visibility, cause, deprecated)
}
