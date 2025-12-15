package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func init() {
	resource.AddTestSweepers("rootly_dashboard", &resource.Sweeper{
		Name: "rootly_dashboard",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListDashboardsParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListDashboardsWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting dashboards, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting dashboards, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting dashboards, got empty response")
				}

				for _, dashboard := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(dashboard.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteDashboardWithResponse(ctx, dashboard.Id)
						if err != nil {
							return fmt.Errorf("Error deleting dashboard: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting dashboard, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted dashboard %s", dashboard.Attributes.Name)
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

func testAccResourceDashboardConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_dashboard" "test" {
	name = "%s"
}
`, name)
}
