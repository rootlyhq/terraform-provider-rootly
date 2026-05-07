package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCommunicationsTemplate(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-comms-tpl")
	typeName := acctest.RandomWithPrefix("tf-comms-type")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsTemplateConfig(rName, typeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_template.test", "name", rName),
				),
			},
		},
	})
}

func testAccResourceCommunicationsTemplateConfig(name, typeName string) string {
	return fmt.Sprintf(`
resource "rootly_communications_type" "test" {
	name = "%s"
}

resource "rootly_communications_template" "test" {
	name                  = "%s"
	communication_type_id = rootly_communications_type.test.id
	body                  = "Incident {{ incident.title }} is {{ incident.status }}"
}
`, typeName, name)
}
