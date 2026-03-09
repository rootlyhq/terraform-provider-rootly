package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCause(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-cause")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCauseConfig(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.test", "name", rName),
					resource.TestCheckResourceAttrSet("rootly_cause.test", "slug"),
				),
			},
			{
				Config: testAccResourceCauseConfig(rName+"-updated", "A description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_cause.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_cause.test", "description", "A description"),
				),
			},
		},
	})
}

func testAccResourceCauseConfig(name, description string) string {
	return fmt.Sprintf(`
resource "rootly_cause" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}
