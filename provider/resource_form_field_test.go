package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormField(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-form-field")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormFieldConfig(rName),
			},
		},
	})
}

func testAccResourceFormFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_field" "test" {
	kind = "custom"
	name = "%s"
}
`, name)
}
