package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendGoogleChatAttachments(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-gc-attach")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendGoogleChatAttachmentsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_send_google_chat_attachments.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_send_google_chat_attachments.foo", "task_params.0.task_type", "send_google_chat_attachments"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskSendGoogleChatAttachmentsConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_send_google_chat_attachments" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
