package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomField(t *testing.T) {
	customFieldLabel := acctest.RandomWithPrefix("tf-custom-field")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomFieldConfig(customFieldLabel),
			},
		},
	})
}

func testAccResourceCustomFieldConfig(customFieldLabel string) string {
	return fmt.Sprintf(`
resource "rootly_custom_field" "testing" {
	label = "%s"
}
`, customFieldLabel)
}
