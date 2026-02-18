package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func init() {
	resource.AddTestSweepers("rootly_team", &resource.Sweeper{
		Name: "rootly_team",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListTeamsParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListTeamsWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting teams, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting teams, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting teams, got empty response")
				}

				for _, team := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(team.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteTeamWithResponse(ctx, team.Id)
						if err != nil {
							return fmt.Errorf("Error deleting team: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting team, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted team %s", team.Attributes.Name)
					}
				}

				if httpResp.ApplicationvndApiJSON200.Links.Next == nil {
					break
				}

				params.PageNumber = ptr.Ptr(ptr.Value(params.PageNumber) + 1)
			}

			return nil
		},
	})
}

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
				Config:            testAccResourceTeamConfig(teamName),
				ConfigStateChecks: configStateChecks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceTeamConfig(teamName),
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
				Config: testAccResourceTeamConfig(teamName),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(teamName)),
				),
			},
			{
				Config: testAccResourceTeamConfig(teamName + "-updated"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(teamName+"-updated")),
				),
			},
		},
	})
}

func testAccResourceTeamConfig(teamName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}
`, teamName)
}
