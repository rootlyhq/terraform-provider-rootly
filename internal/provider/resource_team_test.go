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
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
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
	addr := "rootly_team.test"
	name := acctest.RandomWithPrefix("tf-team")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("slug"), knownvalue.NotNull()),
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
				Config: testAccResourceTeamConfig(name, `
					lifecycle {
						ignore_changes = [user_ids, admin_ids]
					}
				`),
				ConfigStateChecks: configStateChecks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceTeamConfig(name, ""),
				ConfigStateChecks:        configStateChecks,
			},
		},
	})
}

func TestAccResourceTeam(t *testing.T) {
	addr := "rootly_team.test"
	name := acctest.RandomWithPrefix("tf-team")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("slug"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfig(name, ""),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
				),
			},
			{
				Config: testAccResourceTeamConfig(name+"-updated", ""),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
				),
			},
		},
	})
}

func TestAccResourceTeam_WithAlertBroadcast(t *testing.T) {
	addr := "rootly_team.test"
	name := acctest.RandomWithPrefix("tf-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTeamConfig(name, `
					alert_broadcast_enabled = true
					alert_broadcast_channel {
						id   = "id"
						name = "name"
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("alert_broadcast_enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("alert_broadcast_channel"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":   knownvalue.StringExact("id"),
							"name": knownvalue.StringExact("name"),
						}),
					})),
				},
			},
			{
				Config: testAccResourceTeamConfig(name, `
					alert_broadcast_enabled = false
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("alert_broadcast_enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("alert_broadcast_channel"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}

func testAccResourceTeamConfig(name string, extras string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
	%s
}
`, name, extras)
}
