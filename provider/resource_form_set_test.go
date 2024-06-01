package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormSet(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormSet,
			},
		},
	})
}

const testAccResourceFormSet = `
resource "rootly_form_set" "test" {
	name = "test"
forms = ["foo"]
}
`
