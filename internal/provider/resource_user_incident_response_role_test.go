package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccResourceUserIncidentResponseRole_Basic(t *testing.T) {
	addr := "rootly_user_incident_response_role.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserIncidentResponseRoleConfig("user"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("user_id"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(addr, tfjsonpath.New("role_id"), "data.rootly_role.user", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			{
				Config: testAccResourceUserIncidentResponseRoleConfig("observer"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("user_id"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(addr, tfjsonpath.New("role_id"), "data.rootly_role.observer", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			{
				ResourceName: addr,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[addr]
					if !ok {
						return "", fmt.Errorf("not found in state: %s", addr)
					}
					return rs.Primary.Attributes["user_id"], nil
				},
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccResourceUserIncidentResponseRoleConfig(roleName string) string {
	return fmt.Sprintf(`
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}

data "rootly_role" "user" {
	slug = "user"
}

data "rootly_role" "observer" {
	slug = "observer"
}

resource "rootly_user_incident_response_role" "test" {
	user_id = data.rootly_user.test.id
	role_id = data.rootly_role.%[1]s.id
}
`, roleName)
}
