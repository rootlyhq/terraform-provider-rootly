package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceApiKey(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
					resource.TestCheckResourceAttrSet("rootly_api_key.test", "token"),
				),
			},
		},
	})
}

func testAccResourceApiKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "2099-01-01T00:00:00Z"
}
`, name)
}
