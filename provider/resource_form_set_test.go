package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormSet(t *testing.T) {
	formSetName := acctest.RandomWithPrefix("tf-form-set")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormSetConfig(formSetName),
			},
		},
	})
}

func testAccResourceFormSetConfig(formSetName string) string {
	return fmt.Sprintf(`
resource "rootly_form_set" "test" {
	name = "%s"
	forms = ["web_new_incident_form"]
}
`, formSetName)
}
