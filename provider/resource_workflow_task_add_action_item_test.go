package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskAddActionItem(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskAddActionItem,
			},
		},
	})
}

const testAccResourceWorkflowTaskAddActionItem = `
resource "rootly_workflow_incident" "test" {
  name = "test-workflow"
  trigger_params {
    triggers = ["incident_updated"]
  }
}

resource "rootly_workflow_task_add_action_item" "test" {
  workflow_id = rootly_workflow_incident.test.id
  task_params {
    summary  = "Test action item"
    priority = "high"
    status   = "open"
  }
}
`

func TestAccResourceWorkflowTaskAddActionItem_UniquePosition(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowTaskPositionStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.test", "name", "test-position-validation"),
					resource.TestCheckResourceAttr("rootly_workflow_task_print.task1", "position", "1"),
				),
			},
			{
				Config:      testAccWorkflowTaskPositionStep2Duplicate,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`position 1 is already in use`),
			},
			{
				Config: testAccWorkflowTaskPositionStep3Unique,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_print.task1", "position", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_task_print.task2", "position", "2"),
				),
			},
		},
	})
}

const testAccWorkflowTaskPositionStep1 = `
resource "rootly_workflow_incident" "test" {
  name = "test-position-validation"
  trigger_params {
    triggers = ["incident_updated"]
  }
}

resource "rootly_workflow_task_print" "task1" {
  workflow_id = rootly_workflow_incident.test.id
  position    = 1
  task_params {
    message = "First task message"
  }
}
`

const testAccWorkflowTaskPositionStep2Duplicate = `
resource "rootly_workflow_incident" "test" {
  name = "test-position-validation"
  trigger_params {
    triggers = ["incident_updated"]
  }
}

resource "rootly_workflow_task_print" "task1" {
  workflow_id = rootly_workflow_incident.test.id
  position    = 1
  task_params {
    message = "First task message"
  }
}

resource "rootly_workflow_task_print" "task2" {
  workflow_id = rootly_workflow_incident.test.id
  position    = 1  # Duplicate position - should fail validation
  task_params {
    message = "Second task message"
  }
}
`

const testAccWorkflowTaskPositionStep3Unique = `
resource "rootly_workflow_incident" "test" {
  name = "test-position-validation"
  trigger_params {
    triggers = ["incident_updated"]
  }
}

resource "rootly_workflow_task_print" "task1" {
  workflow_id = rootly_workflow_incident.test.id
  position    = 1
  task_params {
    message = "First task message"
  }
}

resource "rootly_workflow_task_print" "task2" {
  workflow_id = rootly_workflow_incident.test.id
  position    = 2  # Unique position - should succeed
  task_params {
    message = "Second task message"
  }
}
`
