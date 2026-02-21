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
	resource.AddTestSweepers("rootly_role", &resource.Sweeper{
		Name: "rootly_role",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListRolesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListRolesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting roles, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting roles, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting roles, got empty response")
				}

				for _, role := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(role.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteRoleWithResponse(ctx, role.Id)
						if err != nil {
							return fmt.Errorf("Error deleting role: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting role, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted role %s", role.Attributes.Name)
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
