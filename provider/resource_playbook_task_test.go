package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePlaybookTask(t *testing.T) {
	playbookName := acctest.RandomWithPrefix("tf-playbook")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePlaybookTaskConfig(playbookName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_playbook_task.test", "task", "Check dashboards"),
				),
			},
		},
	})
}

func testAccResourcePlaybookTaskConfig(playbookName string) string {
	return fmt.Sprintf(`
resource "rootly_playbook" "test" {
	title = "%s"
}

resource "rootly_playbook_task" "test" {
	playbook_id = rootly_playbook.test.id
	task        = "Check dashboards"
}
`, playbookName)
}
