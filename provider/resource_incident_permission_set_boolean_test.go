package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentPermissionSetBoolean(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-incident-permission-set-boolean")

	scopedChecks := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.fully_enabled", "false"),
		resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.applies_to_unassigned", "false"),
		resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.create_enabled", "true"),
		resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.severity_ids.#", "1"),
		resource.TestCheckResourceAttrPair("rootly_incident_permission_set_boolean.test", "severity_params.0.severity_ids.0", "rootly_severity.sev0", "id"),
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, `
	severity_params {
		fully_enabled         = false
		applies_to_unassigned = false
		create_enabled        = true
		severity_ids          = [rootly_severity.sev0.id]
	}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "kind", "update_summary"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "private", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "enabled", "true"),
					scopedChecks,
				),
			},
			// severity_params is a block, so it cannot be cleared with an empty list.
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, `
	severity_params = []`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			// Removing the block keeps the severity scope the API holds.
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, ""),
				Check:  scopedChecks,
			},
			// Resetting the scope takes an explicit block.
			{
				Config: testAccResourceIncidentPermissionSetBooleanConfig(rName, `
	severity_params {
		fully_enabled = true
	}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.fully_enabled", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.applies_to_unassigned", "true"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.create_enabled", "false"),
					resource.TestCheckResourceAttr("rootly_incident_permission_set_boolean.test", "severity_params.0.severity_ids.#", "0"),
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

func testAccResourceIncidentPermissionSetBooleanConfig(name, severityParams string) string {
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
