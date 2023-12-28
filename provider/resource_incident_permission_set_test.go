package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentPermissionSet(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentPermissionSet,
			},
		},
	})
}

const testAccResourceIncidentPermissionSet = `
resource "rootly_incident_permission_set" "test" {
	name = "test"
}
`
