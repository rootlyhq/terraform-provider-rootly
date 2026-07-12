package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowActionItemFormFieldCondition(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-ai-ffc")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowActionItemFormFieldConditionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_action_item_form_field_condition.foo", "action_item_condition", "ANY"),
				),
			},
		},
	})
}

func testAccResourceWorkflowActionItemFormFieldConditionConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_incident" "foo" {
  name = "%s"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_action_item_form_field_condition" "foo" {
  workflow_id          = rootly_workflow_incident.foo.id
  form_field_id        = "placeholder"
  action_item_condition = "ANY"
}
`, name)
}
