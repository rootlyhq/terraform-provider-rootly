package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentPermissionSet(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-permission-set")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentPermissionSetConfig(rName),
			},
		},
	})
}

func testAccResourceIncidentPermissionSetConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_permission_set" "test" {
	name = "%s"
}
`, name)
}
