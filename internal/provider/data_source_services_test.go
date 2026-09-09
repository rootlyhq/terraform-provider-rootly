package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
)

func TestAccDataSourceServices(t *testing.T) {
	rn := "data.rootly_services.test"
	name := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServicesConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("services"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"name": knownvalue.StringExact(name),
						}),
					})),
				},
			},
		},
	})
}

func testAccDataSourceServicesConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_service" "test" {
	name = "%[1]s"

	slack_channels {
		id   = "C08836PQ123"
		name = "terraform"
	}
}

data "rootly_services" "test" {
	depends_on = [resource.rootly_service.test]
}
`, name)
}
