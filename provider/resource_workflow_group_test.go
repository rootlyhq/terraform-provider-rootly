package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkflowGroup(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-workflow-group")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkflowGroupConfig(rName),
			},
		},
	})
}

func testAccResourceWorkflowGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_workflow_group" "test" {
	name = "%s"
}
`, name)
}
