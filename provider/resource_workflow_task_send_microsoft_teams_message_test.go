package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskSendMicrosoftTeamsMessage(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomWithPrefix("tf-wf-task")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			if os.Getenv("ROOTLY_TEST_MS_TEAMS_CHANNEL_ID") == "" {
				t.Skip("ROOTLY_TEST_MS_TEAMS_CHANNEL_ID not set")
			}

			if os.Getenv("ROOTLY_TEST_TEAMS_CHANNEL_NAME") == "" {
				t.Skip("ROOTLY_TEST_TEAMS_CHANNEL_NAME not set")
			}
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskSendMicrosoftTeamsMessageConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_incident.foo", "name", rName),
				),
			},
			{
				Config: testAccResourceWorkflowTaskSendMicrosoftTeamsMessageUpdateConfig(rName),
			},
		},
	})
}

func testAccResourceWorkflowTaskSendMicrosoftTeamsMessageConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
	name = "%[1]s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_microsoft_teams_message" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		text = "test"
		channels {
			id = "%[2]s"
			name = "%[3]s"
		}
	}
}
`, name, os.Getenv("ROOTLY_TEST_MS_TEAMS_CHANNEL_ID"), os.Getenv("ROOTLY_TEST_TEAMS_CHANNEL_NAME"))
}

func testAccResourceWorkflowTaskSendMicrosoftTeamsMessageUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
	name = "%[1]s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}

resource "rootly_workflow_task_send_microsoft_teams_message" "foo" {
	workflow_id = rootly_workflow_incident.foo.id
	task_params {
		text = "test"
		channels {
			id = "%[2]s"
			name = "%[3]s"
		}
	}
}
`, name, os.Getenv("ROOTLY_TEST_MS_TEAMS_CHANNEL_ID"), os.Getenv("ROOTLY_TEST_TEAMS_CHANNEL_NAME"))
}
