package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateSharepointPage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateSharepointPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_sharepoint_page.test", "task_params.0.task_type", "update_sharepoint_page"),
				),
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateSharepointPage = `
resource "rootly_workflow_incident" "test" {
	name = "test-workflow-sharepoint"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_sharepoint_page" "test" {
	workflow_id = rootly_workflow_incident.test.id
	task_params {
		file_id = "01234567-89ab-cdef-0123-456789abcdef"
		title = "Incident Status"
		content = "Updated SharePoint content"
	}
}
`
