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
	resource.AddTestSweepers("rootly_alert_route", &resource.Sweeper{
		Name: "rootly_alert_route",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListAlertRoutesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListAlertRoutesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting alert routes, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting alert routes, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting alert routes, got empty response")
				}

				for _, alertRoute := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(alertRoute.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteAlertRouteWithResponse(ctx, alertRoute.Id)
						if err != nil {
							return fmt.Errorf("Error deleting alert route: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting alert route, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted alert route %s", alertRoute.Attributes.Name)
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
