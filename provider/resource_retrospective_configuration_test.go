package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRetrospectiveConfiguration(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigurationConfig(),
			},
		},
	})
}

func testAccResourceConfigurationConfig() string {
	return `
data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_retrospective_configuration" "test" {
	kind = "skip"
	severity_ids = [data.rootly_severity.critical.id]
}
`
}
