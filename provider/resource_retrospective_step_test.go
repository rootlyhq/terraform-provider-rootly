package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRetrospectiveStep(t *testing.T) {
	retrospectiveProcessName := acctest.RandomWithPrefix("tf-retrospective-process")
	retrospectiveStepTitle := acctest.RandomWithPrefix("tf-retrospective-step")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStepConfig(retrospectiveProcessName, retrospectiveStepTitle),
			},
		},
	})
}

func testAccResourceStepConfig(retrospectiveProcessName, retrospectiveStepTitle string) string {
	return fmt.Sprintf(`
data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_retrospective_process" "test" {
	name = "%s"
	copy_from = "starter_template"
	retrospective_process_matching_criteria {
		severity_ids = [data.rootly_severity.critical.id]
	}
}

resource "rootly_retrospective_step" "test" {
	retrospective_process_id = rootly_retrospective_process.test.id
	title = "%s"
}
`, retrospectiveProcessName, retrospectiveStepTitle)
}
