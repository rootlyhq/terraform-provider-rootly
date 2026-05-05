package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRetrospectiveStep(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-retro-step")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRetrospectiveStepConfig(rName),
			},
		},
	})
}

func testAccResourceRetrospectiveStepConfig(rName string) string {
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
	title = "%s-step"
}
`, rName, rName)
}
