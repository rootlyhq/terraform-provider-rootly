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
	resource.AddTestSweepers("rootly_escalation_policy", &resource.Sweeper{
		Name: "rootly_escalation_policy",
		F: func(region string) error {
			ctx := context.Background()

			params := &rootly.ListEscalationPoliciesParams{
				PageNumber: ptr.Ptr(1),
			}

			for {
				httpResp, err := acctest.SharedClient.ListEscalationPoliciesWithResponse(ctx, params)
				if err != nil {
					return fmt.Errorf("Error getting escalation policies, got error: %s", err)
				} else if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("Error getting escalation policies, got status code: %d", httpResp.StatusCode())
				} else if httpResp.ApplicationvndApiJSON200 == nil {
					return fmt.Errorf("Error getting escalation policies, got empty response")
				}

				for _, escalationPolicy := range httpResp.ApplicationvndApiJSON200.Data {
					if strings.HasPrefix(escalationPolicy.Attributes.Name, "tf-") {
						httpResp, err := acctest.SharedClient.DeleteEscalationPolicyWithResponse(ctx, escalationPolicy.Id)
						if err != nil {
							return fmt.Errorf("Error deleting escalation policy: %s", err)
						} else if httpResp.StatusCode() != http.StatusOK {
							return fmt.Errorf("Error deleting escalation policy, got status code: %d", httpResp.StatusCode())
						}

						log.Printf("[INFO] Deleted escalation policy %s", escalationPolicy.Attributes.Name)
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
