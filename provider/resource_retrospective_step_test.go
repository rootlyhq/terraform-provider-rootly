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
				Config: testAccResourceRetrospectiveStep,
			},
		},
	})
}

const testAccResourceRetrospectiveStep = `
resource "rootly_retrospective_step" "test" {
	title = "test"
}
`
