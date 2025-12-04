// Package polling provides utilities for polling operations until they complete
package polling

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// WaitForAlertRouteRules polls until the alert route has the expected number of rules created
func WaitForAlertRouteRules(ctx context.Context, getRulesFunc func() ([]interface{}, error), expectedCount int) error {
	return retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		rules, err := getRulesFunc()
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error fetching alert route rules: %w", err))
		}

		if len(rules) >= expectedCount {
			return nil // Success - rules are created
		}

		return retry.RetryableError(fmt.Errorf("waiting for alert route rules to be created: expected %d, got %d", expectedCount, len(rules)))
	})
}