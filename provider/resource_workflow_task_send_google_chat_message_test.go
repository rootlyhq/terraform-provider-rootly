package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendGoogleChatMessage(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-gc-msg")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendGoogleChatMessageConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_send_google_chat_message.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_send_google_chat_message.foo", "task_params.0.task_type", "send_google_chat_message"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskSendGoogleChatMessageConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_send_google_chat_message" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
