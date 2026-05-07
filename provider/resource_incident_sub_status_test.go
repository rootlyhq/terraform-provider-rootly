package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIncidentSubStatus(t *testing.T) {
	subStatusName := acctest.RandomWithPrefix("tf-sub")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIncidentSubStatusConfig(subStatusName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_incident_sub_status.test", "id"),
				),
			},
		},
	})
}

func testAccResourceIncidentSubStatusConfig(subStatusName string) string {
	return fmt.Sprintf(`
resource "rootly_sub_status" "test" {
	name          = "%s"
	parent_status = "started"
}

data "rootly_incident" "test" {
	slug = "test"
}

resource "rootly_incident_sub_status" "test" {
	incident_id   = data.rootly_incident.test.id
	sub_status_id = rootly_sub_status.test.id
	assigned_at   = "2026-06-01T00:00:00Z"
}
`, subStatusName)
}
