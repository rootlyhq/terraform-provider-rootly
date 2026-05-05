package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCommunicationsType(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-communications-type")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsTypeConfig(rName),
			},
		},
	})
}

func testAccResourceCommunicationsTypeConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_communications_type" "test" {
	name = "%s"
	color = "#FF0000"
}
`, name)
}
