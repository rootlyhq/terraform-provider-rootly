package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAlertsSource(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-alert-source")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Create resource and query with data source in same step
				Config: testAccDataSourceAlertsSourceConfig(randomName),
				Check: resource.ComposeTestCheckFunc(
					// Verify the resource was created
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "name", randomName),
					// Verify the data source found the correct alert source
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test", "name", randomName),
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test", "source_type", "generic_webhook"),
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test", "status", "setup_incomplete"),
					// Verify the ID matches the resource we created
					resource.TestCheckResourceAttrPair(
						"data.rootly_alerts_source.test", "id",
						"rootly_alerts_source.test", "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceAlertsSource_FilterBySourceType(t *testing.T) {
	randomName := acctest.RandomWithPrefix("tf-test-filter")
	webhookName := fmt.Sprintf("%s-webhook", randomName)
	emailName := fmt.Sprintf("%s-email", randomName)

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resources only, giving the API time to propagate before querying
				Config: testAccDataSourceAlertsSourceFilterResourcesOnlyConfig(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alerts_source.test_webhook", "name", webhookName),
					resource.TestCheckResourceAttr("rootly_alerts_source.test_email", "name", emailName),
				),
			},
			{
				// Step 2: Now query with data sources
				Config: testAccDataSourceAlertsSourceFilterConfig(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test_webhook", "name", webhookName),
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test_webhook", "source_type", "generic_webhook"),
					resource.TestCheckResourceAttrPair(
						"data.rootly_alerts_source.test_webhook", "id",
						"rootly_alerts_source.test_webhook", "id",
					),
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test_email", "name", emailName),
					resource.TestCheckResourceAttr("data.rootly_alerts_source.test_email", "source_type", "email"),
					resource.TestCheckResourceAttrPair(
						"data.rootly_alerts_source.test_email", "id",
						"rootly_alerts_source.test_email", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceAlertsSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
	name        = "%s"
	source_type = "generic_webhook"
}

data "rootly_alerts_source" "test" {
	name        = "%s"
	source_type = "generic_webhook"
	depends_on  = [rootly_alerts_source.test]
}
`, name, name)
}

func testAccDataSourceAlertsSourceFilterResourcesOnlyConfig(name string) string {
	webhookName := fmt.Sprintf("%s-webhook", name)
	emailName := fmt.Sprintf("%s-email", name)

	return fmt.Sprintf(`
resource "rootly_alerts_source" "test_webhook" {
	name        = "%s"
	source_type = "generic_webhook"
}

resource "rootly_alerts_source" "test_email" {
	name        = "%s"
	source_type = "email"
}
`, webhookName, emailName)
}

func testAccDataSourceAlertsSourceFilterConfig(name string) string {
	webhookName := fmt.Sprintf("%s-webhook", name)
	emailName := fmt.Sprintf("%s-email", name)

	return fmt.Sprintf(`
resource "rootly_alerts_source" "test_webhook" {
	name        = "%s"
	source_type = "generic_webhook"
}

resource "rootly_alerts_source" "test_email" {
	name        = "%s"
	source_type = "email"
}

# This should find the generic_webhook source by filtering on name AND source_type
data "rootly_alerts_source" "test_webhook" {
	name        = rootly_alerts_source.test_webhook.name
	source_type = "generic_webhook"
	depends_on  = [rootly_alerts_source.test_webhook, rootly_alerts_source.test_email]
}

# This should find the email source by filtering on name AND source_type
data "rootly_alerts_source" "test_email" {
	name        = rootly_alerts_source.test_email.name
	source_type = "email"
	depends_on  = [rootly_alerts_source.test_webhook, rootly_alerts_source.test_email]
}
`, webhookName, emailName)
}
