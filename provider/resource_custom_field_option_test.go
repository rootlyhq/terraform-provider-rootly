package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCustomFieldOption(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-cfo")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCustomFieldOptionConfig(rName, rName+"-opt"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field.parent", "label", rName),
					resource.TestCheckResourceAttr("rootly_custom_field_option.foo", "value", rName+"-opt"),
				),
			},
			{
				Config: testAccResourceCustomFieldOptionConfig(rName, rName+"-opt2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_custom_field_option.foo", "value", rName+"-opt2"),
				),
			},
		},
	})
}

func testAccResourceCustomFieldOptionConfig(label, optionValue string) string {
	return fmt.Sprintf(`
resource "rootly_custom_field" "parent" {
  label = "%s"
	kind = "select"
	shown = ["incident_form", "incident_slack_form"]
	required = []
}

resource "rootly_custom_field_option" "foo" {
	custom_field_id = rootly_custom_field.parent.id
  value = "%s"
}
`, label, optionValue)
}
