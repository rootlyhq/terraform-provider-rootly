package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Once tasks exist at positions 1-5, moving the first task to the
// end (and renumbering the rest) should apply cleanly.
func TestAccWorkflowTaskPositionReorderDeadlock(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-wf-task-reorder")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowTaskPositionReorderConfig(rName, []int{1, 2, 3, 4, 5}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.a", "position", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.b", "position", "2"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.c", "position", "3"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.d", "position", "4"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.e", "position", "5"),
				),
			},
			{
				Config: testAccWorkflowTaskPositionReorderConfig(rName, []int{5, 1, 2, 3, 4}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.a", "position", "5"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.b", "position", "1"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.c", "position", "2"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.d", "position", "3"),
					resource.TestCheckResourceAttr("rootly_workflow_task_add_to_timeline.e", "position", "4"),
				),
			},
		},
	})
}

func testAccWorkflowTaskPositionReorderConfig(name string, positions []int) string {
	var b strings.Builder
	fmt.Fprintf(&b, `
resource "rootly_workflow_incident" "foo" {
	name = "%s"
	trigger_params {
		triggers = ["incident_updated"]
	}
}
`, name)

	// Chain creates with depends_on so tasks are created one at a time.
	prev := ""
	for i, position := range positions {
		key := string(rune('a' + i))
		dependsOn := ""
		if prev != "" {
			dependsOn = fmt.Sprintf("depends_on = [rootly_workflow_task_add_to_timeline.%s]", prev)
		}
		fmt.Fprintf(&b, `
resource "rootly_workflow_task_add_to_timeline" "%s" {
	workflow_id = rootly_workflow_incident.foo.id
	name        = "Task %s"
	position    = %d
	task_params {
		event = "task %s"
	}
	%s
}
`, key, strings.ToUpper(key), position, key, dependsOn)
		prev = key
	}
	return b.String()
}
