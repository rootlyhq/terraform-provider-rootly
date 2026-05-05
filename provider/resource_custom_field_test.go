package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomField(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-custom-field")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomFieldConfig(rLabel),
			},
		},
	})
}

func testAccResourceCustomFieldConfig(label string) string {
	return fmt.Sprintf(`
resource "rootly_custom_field" "test" {
	label = "%s"
}
`, label)
}
