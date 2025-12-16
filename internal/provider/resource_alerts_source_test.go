package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	resource.AddTestSweepers("rootly_alerts_source", &resource.Sweeper{
		Name: "rootly_alerts_source",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListAlertsSourcesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListAlertsSourcesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting alerts sources, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting alerts sources, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting alerts sources, got empty response")
				}

				for _, alertSource := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(alertSource.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteAlertsSourceWithResponse(ctx, alertSource.Id)
						if err != nil {
							return fmt.Errorf("Error deleting alerts source: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting alerts source, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted alerts source %s", alertSource.Attributes.Name)
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

func TestAccResourceAlertsSource(t *testing.T) {
	resName := "rootly_alerts_source.test"
	teamName := acctest.RandomWithPrefix("tf-team")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfig(teamName, alertUrgencyName, alertsSourceName, `
					source_type = "generic_webhook"

					deduplicate_alerts_by_key = true
					deduplication_key_path = "$.id"

					alert_source_urgency_rules_attributes {
						alert_urgency_id = rootly_alert_urgency.test.id
						json_path = "test"
						operator = "is"
						value = "P1"
					}

					sourceable_attributes {
						auto_resolve  = false
						resolve_state = "$.status"
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(alertsSourceName)),
				),
			},
			{
				Config: testAccResourceAlertsSourceConfig(teamName, alertUrgencyName, alertsSourceName+"-updated", `
					source_type = "generic_webhook"
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(alertsSourceName+"-updated")),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_AlertTemplateAttributesErrorWhenAlertFieldsEnabled(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfig(teamName, alertUrgencyName, alertsSourceName, `
					source_type = "generic_webhook"

					alert_template_attributes {
						title = "alert title"
						description = "alert description"
						external_url = "https://example.com"
					}
				`),
				ExpectError: regexp.MustCompile("Alert template attributes cannot be provided when alert fields are enabled"),
			},
		},
	})
}

func testAccResourceAlertsSourceConfig(teamName, alertUrgencyName, alertsSourceName, extra string) string {
	return testAccResourceTeamConfig(teamName) + testAccResourceAlertUrgencyConfig(alertUrgencyName, alertUrgencyName+" description") + fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
	depends_on = [rootly_alert_urgency.test]

	name            = "%s"
	owner_group_ids = [rootly_team.test.id]

	%s
}
`, alertsSourceName, extra)
}
