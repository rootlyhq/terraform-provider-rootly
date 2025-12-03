package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
)

func TestAccResourceTeam_UpgradeFromVersion(t *testing.T) {
	resName := "rootly_team.test"
	teamName := acctest.RandomWithPrefix("tf-team")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(teamName)),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("slug"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rootly": {
						Source:            "rootlyhq/rootly",
						VersionConstraint: "4.3.8",
					},
				},
				Config:            testAccResourceTeam(teamName),
				ConfigStateChecks: configStateChecks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceTeam(teamName),
				ConfigStateChecks:        configStateChecks,
			},
		},
	})
}

func TestAccResourceTeam(t *testing.T) {
	resName := "rootly_team.test"
	teamName := acctest.RandomWithPrefix("tf-team")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("slug"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeam(teamName),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(teamName)),
				),
			},
			{
				Config: testAccResourceTeam(teamName + "-updated"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(teamName+"-updated")),
				),
			},
		},
	})
}

func testAccResourceTeam(teamName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}
`, teamName)
}
