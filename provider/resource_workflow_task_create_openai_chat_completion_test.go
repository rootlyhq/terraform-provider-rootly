package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowTaskCreateOpenaiChatCompletion(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-openai")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowTaskCreateOpenaiChatCompletionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_create_openai_chat_completion.foo", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_workflow_task_create_openai_chat_completion.foo", "task_params.0.task_type", "create_openai_chat_completion"),
				),
			},
		},
	})
}

func testAccResourceWorkflowTaskCreateOpenaiChatCompletionConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_create_openai_chat_completion" "foo" {
  workflow_id = rootly_workflow_incident.foo.id
  task_params {
    model = {
      id   = "foo"
      name = "bar"
    }
    prompt     = "test prompt"
    max_tokens = 100
  }
}
`, name)
}
