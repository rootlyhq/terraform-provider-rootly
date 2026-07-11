package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskChangeGoogleChatSpacePrivacy(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-gc-priv")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskChangeGoogleChatSpacePrivacyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_change_google_chat_space_privacy.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_change_google_chat_space_privacy.foo", "task_params.0.task_type", "change_google_chat_space_privacy"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskChangeGoogleChatSpacePrivacyConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_change_google_chat_space_privacy" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
