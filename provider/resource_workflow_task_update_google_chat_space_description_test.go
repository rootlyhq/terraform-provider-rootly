package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskUpdateGoogleChatSpaceDescription(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-gc-desc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskUpdateGoogleChatSpaceDescriptionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_update_google_chat_space_description.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_update_google_chat_space_description.foo", "task_params.0.task_type", "update_google_chat_space_description"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskUpdateGoogleChatSpaceDescriptionConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_update_google_chat_space_description" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
