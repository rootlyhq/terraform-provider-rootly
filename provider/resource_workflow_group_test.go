package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_group.foo", "name", "Default"),
				),
			},
			{
				Config: testAccResourceWorkflowGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_group.foo", "name", "Notifications"),
				),
			},
		},
	})
}

const testAccResourceWorkflowGroup = `
resource "rootly_workflow_group" "foo" {
  name = "Default"
	kind = "incident"
}
`

const testAccResourceWorkflowGroupUpdate = `
resource "rootly_workflow_group" "foo" {
  name = "Notifications"
	kind = "incident"
}
`
