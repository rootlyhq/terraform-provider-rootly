package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentPermissionSetBoolean(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-permission-set-boolean")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "kind", "update_summary"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "private", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.fully_enabled", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.applies_to_unassigned", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.create_enabled", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.severity_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_incident_permission_set_boolean.test", "severity_params.0.severity_ids.0", "rootly_severity.sev0", "id"),
				),
			},
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.fully_enabled", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.applies_to_unassigned", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.create_enabled", "false"),
				),
			},
			{
				ResourceName:      "rootly_incident_permission_set_boolean.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceIncidentPermissionSetBooleanConfig(name string, fullyEnabled bool) string {
	severityParams := `
	severity_params {
		fully_enabled         = false
		applies_to_unassigned = false
		create_enabled        = true
		severity_ids          = [rootly_severity.sev0.id]
	}`
	if fullyEnabled {
		severityParams = `
	severity_params {
		fully_enabled = true
	}`
	}

	return fmt.Sprintf(`
resource "rootly_severity" "sev0" {
	name = "%[1]s-sev0"
}

resource "rootly_incident_permission_set" "test" {
	name                         = "%[1]s"
	description                  = "%[1]s-description"
	private_incident_permissions = ["read", "update"]
}

resource "rootly_incident_permission_set_boolean" "test" {
	incident_permission_set_id = rootly_incident_permission_set.test.id
	kind                       = "update_summary"
	private                    = true
	enabled                    = true
%[2]s
}
`, name, severityParams)
}
