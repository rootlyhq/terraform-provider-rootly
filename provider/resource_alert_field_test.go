package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertField(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-alert-field")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertFieldConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", rName),
				),
			},
			{
				Config: testAccResourceAlertFieldConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", rName+"-updated"),
				),
			},
		},
	})
}

func testAccResourceAlertFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_field" "test" {
  name = "%s"
}
`, name)
}
