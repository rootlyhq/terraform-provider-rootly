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
