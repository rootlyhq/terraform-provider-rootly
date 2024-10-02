package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSubStatus(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSubStatus,
			},
		},
	})
}

const testAccResourceSubStatus = `
resource "rootly_sub_status" "test" {
	name = "test"
	parent_status = "started"
}
`
