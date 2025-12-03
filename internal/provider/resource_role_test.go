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

func TestAccResourceRole_UpgradeFromVersion(t *testing.T) {
	resName := "rootly_role.test"
	roleName := acctest.RandomWithPrefix("tf-role")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(roleName)),
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
				Config:            testAccResourceRole(roleName),
				ConfigStateChecks: configStateChecks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceRole(roleName),
				ConfigStateChecks:        configStateChecks,
			},
		},
	})
}

func TestAccResourceRole(t *testing.T) {
	resName := "rootly_role.test"
	roleName := acctest.RandomWithPrefix("tf-role")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("slug"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRole(roleName),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(roleName)),
				),
			},
			{
				Config: testAccResourceRole(roleName + "-updated"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(roleName+"-updated")),
				),
			},
		},
	})
}

func testAccResourceRole(roleName string) string {
	return fmt.Sprintf(`
resource "rootly_role" "test" {
	name = "%s"
}
`, roleName)
}
