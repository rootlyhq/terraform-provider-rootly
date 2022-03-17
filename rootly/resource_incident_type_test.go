package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "name", "myincidenttype"),
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "description", ""),
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "color", "#047BF8"),
				),
			},
			{
				Config: testAccResourceIncidentTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "name", "myincidenttype2"),
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_incident_type.foo", "color", "#203"),
				),
			},
		},
	})
}

const testAccResourceIncidentType = `
resource "rootly_incident_type" "foo" {
  name = "myincidenttype"
}
`

const testAccResourceIncidentTypeUpdate = `
resource "rootly_incident_type" "foo" {
  name        = "myincidenttype2"
  description = "test description"
  color       = "#203"
}
`
