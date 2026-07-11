package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskArchiveGoogleChatSpaces(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-gc-archive")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskArchiveGoogleChatSpacesConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_archive_google_chat_spaces.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_archive_google_chat_spaces.foo", "task_params.0.task_type", "archive_google_chat_spaces"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskArchiveGoogleChatSpacesConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_archive_google_chat_spaces" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
