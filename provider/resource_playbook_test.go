package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePlaybook(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-playbook")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePlaybookConfig(rName),
			},
		},
	})
}

func testAccResourcePlaybookConfig(title string) string {
	return fmt.Sprintf(`
resource "rootly_playbook" "test" {
	title = "%s"
}
`, title)
}
