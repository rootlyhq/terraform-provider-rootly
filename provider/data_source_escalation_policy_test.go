package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEscalationPolicy_byName(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ds-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEscalationPolicyByNameConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rootly_escalation_policy.by_name", "name", rName),
					resource.TestCheckResourceAttrPair(
						"data.rootly_escalation_policy.by_name", "id",
						"rootly_escalation_policy.test", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceEscalationPolicyByNameConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_escalation_policy" "test" {
  name = "%s"
}

data "rootly_escalation_policy" "by_name" {
  name       = rootly_escalation_policy.test.name
  depends_on = [rootly_escalation_policy.test]
}
`, name)
}
