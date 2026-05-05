package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentRole(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-role")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentRoleConfig(rName),
			},
		},
	})
}

func testAccResourceIncidentRoleConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "test" {
	name = "%s"
}
`, name)
}
