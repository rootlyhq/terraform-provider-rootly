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
	resource.AddTestSweepers("rootly_alert_urgency", &resource.Sweeper{
		Name: "rootly_alert_urgency",
		Dependencies: []string{
			"rootly_alerts_source",
		},
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListAlertUrgenciesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListAlertUrgenciesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting alert urgencies, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting alert urgencies, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting alert urgencies, got empty response")
				}

				for _, alertUrgency := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(alertUrgency.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteAlertUrgencyWithResponse(ctx, alertUrgency.Id)
						if err != nil {
							return fmt.Errorf("Error deleting alert urgency: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting alert urgency, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted alert urgency %s", alertUrgency.Attributes.Name)
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

func TestAccResourceAlertUrgency(t *testing.T) {
	resName := "rootly_alert_urgency.test"
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("position"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertUrgency(alertUrgencyName, "description"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(alertUrgencyName)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("description")),
				),
			},
			{
				Config: testAccResourceAlertUrgency(alertUrgencyName+"-updated", "updated description"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(alertUrgencyName+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("updated description")),
				),
			},
		},
	})
}

func testAccResourceAlertUrgency(alertUrgencyName, alertUrgencyDescription string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name = "%s"
	description = "%s"
}
`, alertUrgencyName, alertUrgencyDescription)
}
