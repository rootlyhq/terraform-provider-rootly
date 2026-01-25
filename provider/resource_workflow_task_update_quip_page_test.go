package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateQuipPage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateQuipPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_quip_page.test", "task_params.0.task_type", "update_quip_page"),
				),
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateQuipPage = `
resource "rootly_workflow_incident" "test" {
	name = "test-workflow-quip"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_quip_page" "test" {
	workflow_id = rootly_workflow_incident.test.id
	task_params {
		file_id = "AbCdEfGhIj12"
		title = "Incident Documentation"
		content = "Updated Quip content"
	}
}
`
