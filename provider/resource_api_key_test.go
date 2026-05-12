package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceApiKey(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyConfig(rName, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
				),
			},
		},
	})
}

func testAccResourceApiKeyConfig(name, expiresAt string) string {
	return fmt.Sprintf(`
resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "%s"
}
`, name, expiresAt)
}
