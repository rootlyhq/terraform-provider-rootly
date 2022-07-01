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
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", "test-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowAlertUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "name", "test-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_alert.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowAlert = `
resource "rootly_workflow_alert" "foo" {
  name = "test-workflow"
	trigger_params {
		triggers = ["alert_created"]
	}
}
`

const testAccResourceWorkflowAlertUpdate = `
resource "rootly_workflow_alert" "foo" {
  name       = "test-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["alert_created"]
	}
}
`
