package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateConfluencePage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateConfluencePage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_confluence_page.test", "task_params.0.task_type", "update_confluence_page"),
				),
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateConfluencePage = `
resource "rootly_workflow_incident" "test" {
	name = "test-workflow-confluence"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_confluence_page" "test" {
	workflow_id = rootly_workflow_incident.test.id
	task_params {
		file_id = "123456789"
		title = "Incident Update"
		content = "Updated incident details"
	}
}
`
