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

func TestAccDataSourceService(t *testing.T) {
	rn := "data.rootly_service.test"
	name := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slack_channels"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":   knownvalue.StringExact("C08836PQ123"),
							"name": knownvalue.StringExact("terraform"),
						}),
					})),
				},
			},
		},
	})
}

func testAccDataSourceServiceConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_service" "test" {
	name = "%[1]s"

	slack_channels {
		id   = "C08836PQ123"
		name = "terraform"
	}
}

data "rootly_service" "test" {
	id = resource.rootly_service.test.id
}
`, name)
}
