package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRetrospectiveStep(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStep,
			},
		},
	})
}

const testAccResourceStep = `
data "rootly_severity" "critical" {
  slug = "sev0-f3b589d9-732e-4a42-aab5-d1c590c72362"
}

resource "rootly_retrospective_process" "test" {
	name = "testing-tf"
	copy_from = "starter_template"
	retrospective_process_matching_criteria {
		severity_ids = [data.rootly_severity.critical.id]
	}
}

resource "rootly_retrospective_step" "test" {
	retrospective_process_id = rootly_retrospective_process.test.id
	title = "testing-tf"
}
`
