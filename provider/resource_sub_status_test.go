package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSubStatus(t *testing.T) {
	subStatusName := acctest.RandomWithPrefix("tf-sub-status")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSubStatusConfig(subStatusName),
			},
		},
	})
}

func testAccResourceSubStatusConfig(subStatusName string) string {
	return fmt.Sprintf(`
resource "rootly_sub_status" "test" {
	name = "%s"
	parent_status = "started"
}
`, subStatusName)
}
