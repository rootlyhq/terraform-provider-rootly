package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOnCallRole(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-on-call-role")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOnCallRoleConfig(rName),
			},
		},
	})
}

func testAccResourceOnCallRoleConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_on_call_role" "test" {
	name = "%s"
}
`, name)
}

func testAccResourceOnCallRoleWithPermissionsConfig(roleName, roleDisplayName string) string {
	return fmt.Sprintf(`
resource "rootly_on_call_role" "%s" {
	name = "%s"
	alert_fields_permissions = ["create", "read", "update", "delete"]
	alert_groups_permissions = ["create", "read", "update", "delete"]
	alert_routing_rules_permissions = ["create", "read", "update", "delete"]
	on_call_readiness_report_permissions = ["read"]
	on_call_roles_permissions = ["create", "read", "update", "delete"]
	schedule_override_permissions = ["create", "update"]
}
`, roleName, roleDisplayName)
}

func TestAccResourceOnCallRoleWithPermissions(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOnCallRoleWithPermissionsConfig(acctest.RandomWithPrefix("tf-on-call-role"), acctest.RandomWithPrefix("tf-role-display-name")),
			},
		},
	})
}
