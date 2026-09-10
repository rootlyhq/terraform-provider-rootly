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

func TestAccResourceService(t *testing.T) {
	addr := "rootly_service.test"
	name := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceConfig(name, `
					slack_aliases {
						id   = "S0883KV6123"
						name = "eng-terraform"
					}
					slack_channels {
						id   = "C08836PQ123"
						name = "terraform"
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_aliases"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":   knownvalue.StringExact("S0883KV6123"),
							"name": knownvalue.StringExact("eng-terraform"),
						}),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_channels"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":   knownvalue.StringExact("C08836PQ123"),
							"name": knownvalue.StringExact("terraform"),
						}),
					})),
				},
			},
			{
				Config: testAccResourceServiceConfig(name+"-updated", `
					slack_channels {
						id   = "C08836PQ123"
						name = "terraform"
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_aliases"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_channels"), knownvalue.ListExact([]knownvalue.Check{
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

func testAccResourceServiceConfig(name, extra string) string {
	return fmt.Sprintf(`
resource "rootly_service" "test" {
	name = "%[1]s"
	%[2]s
}
`, name, extra)
}
