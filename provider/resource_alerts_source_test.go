package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertsSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSource,
			},
		},
	})
}

const testAccResourceAlertsSource = `
resource "rootly_alerts_source" "test" {
  name = "Generic!!!!!!!!!!"
  source_type = "generic_webhook"
  alert_urgency_id = "066ee186-fc93-45ce-a316-c2a9ac788880"

  alert_template_attributes {
    title = "Server exploded"
    description = "Datacenter is burning down."
    external_url = "https://rootly.com"
  }

  sourceable_attributes {
    auto_resolve = true
    resolve_state = "$.status"

    field_mappings_attributes {
      field = "state"
      json_path = "$.my_group_attribute"
    }

    field_mappings_attributes {
      field = "external_id"
      json_path = "$.my_id_attribute"
    }
  }
}
`

