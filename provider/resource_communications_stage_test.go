package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCommunicationsStage(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-communications-stage")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsStageConfig(rName),
			},
		},
	})
}

func testAccResourceCommunicationsStageConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_communications_stage" "test" {
	name = "%s"
}
`, name)
}
