package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentType(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-type")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentTypeConfig(rName),
			},
		},
	})
}

func testAccResourceIncidentTypeConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_type" "test" {
	name = "%s"
}
`, name)
}
