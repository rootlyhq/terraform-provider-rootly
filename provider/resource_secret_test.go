package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecret(t *testing.T) {

	// We can't use the acctest.RandomWithPrefix() because
	// it will generate a string output that contains a hyphen which breaks
	// our validation rules - so we use the underlying SDK's RandInt() via math/rand
	secretName := fmt.Sprintf("tf_secret_%d", rand.Intn(10000))
	secretValue := fmt.Sprintf("tf_secret_value_%d", rand.Intn(10000))
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretConfig(secretName, secretValue),
			},
		},
	})
}

func testAccResourceSecretConfig(secretName, secretValue string) string {
	return fmt.Sprintf(`
resource "rootly_secret" "test" {
	name = "%s"
	secret = "%s"
}
`, secretName, secretValue)
}
