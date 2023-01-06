package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowGroup,
			},
		},
	})
}

const testAccResourceWorkflowGroup = `
resource "rootly_workflow_group" "test" {
	kind = "simple"
name = "test"
}
`
