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

func TestAccDataSourceService_Validation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_Empty,
				ExpectError: acctest.ExpectLiteralErrors(`At least one attribute out of [id,name,slug,backstage_id,cortex_id,external_id,alert_broadcast_enabled,incident_broadcast_enabled] must be specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "name" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "slug" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "backstage_id" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "cortex_id" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "external_id" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "alert_broadcast_enabled" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly:    true,
				Config:      testAccDataSourceServiceConfig_InvalidFields,
				ExpectError: acctest.ExpectLiteralErrors(`Attribute "incident_broadcast_enabled" cannot be specified when "id" is specified`),
			},
		},
	})
}

func TestAccDataSourceService_ById(t *testing.T) {
	rn := "data.rootly_service.test"
	name := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceConfig(name, `
					id = resource.rootly_service.test.id
				`),
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

func TestAccDataSourceService_ByName(t *testing.T) {
	rn := "data.rootly_service.test"
	name := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceConfig(name, `
					name = resource.rootly_service.test.name
				`),
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

const testAccDataSourceServiceConfig_Empty = `
data "rootly_service" "test" {
}
`

const testAccDataSourceServiceConfig_InvalidFields = `
data "rootly_service" "test" {
	id = "123"
	name = "name"
	slug = "slug"
	backstage_id = "123"
	cortex_id = "123"
	external_id = "123"
	alert_broadcast_enabled = true
	incident_broadcast_enabled = true
}
`

func testAccDataSourceServiceConfig(name, extras string) string {
	return fmt.Sprintf(`
resource "rootly_service" "test" {
	name = "%[1]s"

	slack_channels {
		id   = "C08836PQ123"
		name = "terraform"
	}
}

data "rootly_service" "test" {
	%[2]s
}
`, name, extras)
}
