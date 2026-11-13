package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentPermissionSetResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-permission-set-resource")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentPermissionSetResourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "kind", "severities"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "private", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "resource_type", "Severity"),
					resource.TestCheckResourceAttrPair("rootly_incident_permission_set_resource.test", "resource_id", "rootly_severity.sev1", "id"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "severity_params.0.fully_enabled", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "severity_params.0.applies_to_unassigned", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "severity_params.0.create_enabled", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_resource.test", "severity_params.0.severity_ids.#", "2"),
				),
			},
			{
				ResourceName:      "rootly_incident_permission_set_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceIncidentPermissionSetResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_severity" "sev0" {
	name = "%[1]s-sev0"
}

resource "rootly_severity" "sev1" {
	name = "%[1]s-sev1"
}

resource "rootly_incident_permission_set" "test" {
	name                        = "%[1]s"
	description                 = "%[1]s-description"
	public_incident_permissions = ["read", "update"]
}

resource "rootly_incident_permission_set_resource" "test" {
	incident_permission_set_id = rootly_incident_permission_set.test.id
	kind                       = "severities"
	private                    = false
	resource_type              = "Severity"
	resource_id                = rootly_severity.sev1.id

	severity_params {
		fully_enabled = false
		severity_ids  = [rootly_severity.sev0.id, rootly_severity.sev1.id]
	}
}
`, name)
}
