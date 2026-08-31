package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertUrgency(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-urgency")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertUrgencyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "description", "Test urgency"),
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "retrigger_timeout_minutes", "20"),
				),
			},
			{
				Config: testAccResourceAlertUrgencyConfigUpdated(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "description", "Updated urgency"),
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "retrigger_timeout_minutes", "-1"),
				),
			},
			{
				// Clearing retrigger_timeout_minutes must round-trip to null (inherit workspace default).
				Config: testAccResourceAlertUrgencyConfigCleared(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_urgency.test", "retrigger_timeout_minutes", "0"),
				),
			},
		},
	})
}

func testAccResourceAlertUrgencyConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name        = "%s"
	description = "Test urgency"
	retrigger_timeout_minutes = 20
}
`, name)
}

func testAccResourceAlertUrgencyConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name        = "%s"
	description = "Updated urgency"
	retrigger_timeout_minutes = -1
}
`, name)
}

func testAccResourceAlertUrgencyConfigCleared(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name        = "%s"
	description = "Updated urgency"
}
`, name)
}
