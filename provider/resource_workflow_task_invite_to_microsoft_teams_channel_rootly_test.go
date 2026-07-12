package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskInviteToMicrosoftTeamsChannelRootly(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-mst-inv")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskInviteToMicrosoftTeamsChannelRootlyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_invite_to_microsoft_teams_channel_rootly.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_invite_to_microsoft_teams_channel_rootly.foo", "task_params.0.task_type", "invite_to_microsoft_teams_channel_rootly"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskInviteToMicrosoftTeamsChannelRootlyConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_invite_to_microsoft_teams_channel_rootly" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
  }
}
`, name)
}
