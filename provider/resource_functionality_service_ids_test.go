package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFunctionality_serviceIdsWithNewService(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("tf-service")
	functionalityName := acctest.RandomWithPrefix("tf-functionality")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFunctionalityServiceIdsConfig(serviceName, functionalityName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_functionality.test", "name", functionalityName),
					resource.TestCheckResourceAttr("rootly_functionality.test", "service_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_functionality.test", "service_ids.0", "rootly_service.test", "id"),
				),
			},
			{
				Config:   testAccResourceFunctionalityServiceIdsConfig(serviceName, functionalityName),
				PlanOnly: true,
			},
		},
	})
}

func testAccResourceFunctionalityServiceIdsConfig(serviceName, functionalityName string) string {
	return fmt.Sprintf(`
resource "rootly_service" "test" {
  name = "%s"
}

resource "rootly_functionality" "test" {
  name       = "%s"
  service_ids = [rootly_service.test.id]
}
`, serviceName, functionalityName)
}
