package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowIncident(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowIncident,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-incident-workflow"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceWorkflowIncidentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", "test-incident-workflow2"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceWorkflowIncident = `
resource "rootly_workflow_incident" "foo" {
  name = "test-incident-workflow"
	trigger_params {
		triggers = ["incident_updated"]
	}
}
`

const testAccResourceWorkflowIncidentUpdate = `
resource "rootly_workflow_incident" "foo" {
  name       = "test-incident-workflow2"
  description = "test description"
  enabled     = false
	trigger_params {
		triggers = ["incident_updated", "incident_created"]
	}
}
`
