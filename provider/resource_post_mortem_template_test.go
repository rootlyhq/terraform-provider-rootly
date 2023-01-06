package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePostmortemTemplate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePostmortemTemplate,
			},
		},
	})
}

const testAccResourcePostmortemTemplate = `
resource "rootly_postmortem_template" "test" {
	name = "test"
	content = "test"
}
`
