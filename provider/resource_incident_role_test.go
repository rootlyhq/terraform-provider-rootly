package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentRole(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
			time.Sleep(1 * time.Second)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "name", "myincidentrole"),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "summary", ""),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceIncidentRoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "name", "myincidentrole2"),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "summary", "my summary"),
					resource.TestCheckResourceAttr("rootly_incident_role.foo", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceIncidentRole = `
resource "rootly_incident_role" "foo" {
  name = "myincidentrole"
}
`

const testAccResourceIncidentRoleUpdate = `
resource "rootly_incident_role" "foo" {
  name        = "myincidentrole2"
  description = "test description"
  summary     = "my summary"
  enabled     = false
}
`
