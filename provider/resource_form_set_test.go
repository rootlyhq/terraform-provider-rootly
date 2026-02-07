package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormSet(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-form-set")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormSet(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_set.test", "name", randomName),
				),
			},
		},
	})
}

func testAccResourceFormSet(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_set" "test" {
	name = "%s"
	forms = ["web_new_incident_form"]
}
`, name)
}
