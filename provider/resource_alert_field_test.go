package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertField(t *testing.T) {
	alertFieldName := acctest.RandomWithPrefix("tf-alert-field")
	alertFieldNameUpdated := acctest.RandomWithPrefix("tf-alert-field-updated")
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertFieldConfig(alertFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", alertFieldName),
				),
			},
			{
				Config: testAccResourceAlertFieldConfig(alertFieldNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_field.test", "name", alertFieldNameUpdated),
				),
			},
		},
	})
}

func testAccResourceAlertFieldConfig(alertFieldName string) string {
	return fmt.Sprintf(`
resource "rootly_alert_field" "test" {
  name = "%s"
}
`, alertFieldName)
}
