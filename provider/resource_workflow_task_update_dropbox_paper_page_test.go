package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateDropboxPaperPage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateDropboxPaperPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_dropbox_paper_page.test", "task_params.0.task_type", "update_dropbox_paper_page"),
				),
			},
		},
	})
}

const testAccResourceWorkflowTaskUpdateDropboxPaperPage = `
resource "rootly_workflow_incident" "test" {
	name = "test-workflow-dropbox"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_update_dropbox_paper_page" "test" {
	workflow_id = rootly_workflow_incident.test.id
	task_params {
		file_id = "dPW7HxvCTb6TGr5F2funM"
		title = "Incident Report"
		content = "Updated paper content"
	}
}
`
