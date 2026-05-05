package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecret(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-secret")
	rSecret := acctest.RandomWithPrefix("tf-secret-value")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretConfig(rName, rSecret),
			},
		},
	})
}

func testAccResourceSecretConfig(name, secret string) string {
	return fmt.Sprintf(`
resource "rootly_secret" "test" {
	name = "%s"
	secret = "%s"
}
`, name, secret)
}
