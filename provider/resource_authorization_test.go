package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAuthorization(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthorization,
			},
		},
	})
}

const testAccResourceAuthorization = `
resource "rootly_authorization" "test" {
	authorizable_id = "test"
authorizable_type = "Dashboard"
grantee_id = "test"
grantee_type = "User"
permissions = ["foo"]
}
`
