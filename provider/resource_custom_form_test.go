package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomForm(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-custom-form")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomFormConfig(rName),
			},
		},
	})
}

func testAccResourceCustomFormConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_custom_form" "test" {
	name = "%s"
	command = "test"
}
`, name)
}
