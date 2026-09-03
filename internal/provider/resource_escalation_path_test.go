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

func TestAccResourceEscalationPath_UpgradeFromVersion(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rootly": {
						Source:            "rootlyhq/rootly",
						VersionConstraint: "5.20.1",
					},
				},
				Config: testAccResourceEscalationPathConfig(name, ``),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceEscalationPathConfig(name, ``),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccResourceEscalationPath(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPathConfig(name, `
					initial_delay = 5
					time_restriction_time_zone = "America/New_York"
					time_restrictions {
						start_day = "monday"
						start_time = "17:00"
						end_day = "tuesday"
						end_time = "07:00"
					}
					time_restrictions {
						start_day = "tuesday"
						start_time = "17:00"
						end_day = "wednesday"
						end_time = "07:00"
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("initial_delay"), knownvalue.Int64Exact(5)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restriction_time_zone"), knownvalue.StringExact("America/New_York")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restrictions"), knownvalue.SetSizeExact(2)),
				},
			},
			{
				Config: testAccResourceEscalationPathConfig(name, `
					initial_delay = 0
					time_restriction_time_zone = "Europe/London"
					time_restrictions {
						start_day = "monday"
						start_time = "17:00"
						end_day = "tuesday"
						end_time = "07:00"
					}

				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("initial_delay"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restriction_time_zone"), knownvalue.StringExact("Europe/London")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restrictions"), knownvalue.SetSizeExact(1)),
				},
			},
		},
	})
}

func testAccResourceEscalationPathConfig(name, extra string) string {
	return fmt.Sprintf(`
resource "rootly_escalation_policy" "test" {
	name = "%[1]s-ep"
}

resource "rootly_escalation_path" "test" {
	escalation_policy_id = rootly_escalation_policy.test.id

	name = "%[1]s-path"
	%[2]s
}
`, name, extra)
}
