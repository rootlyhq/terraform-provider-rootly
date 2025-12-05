package polling

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type AsyncRuleCreationStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type AsyncRuleCreationStatusFunc func(alertRouteID, requestId string) (*AsyncRuleCreationStatus, error)

type RefetchAlertRouteFunc func() error

func WaitForAsyncRuleCreation(ctx context.Context, alertRouteID, requestId string, checkStatusFunc AsyncRuleCreationStatusFunc) error {
	return retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		status, err := checkStatusFunc(alertRouteID, requestId)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error checking async rule creation status: %w", err))
		}

		switch status.Status {
		case "success":
			return nil
		case "pending":
			return retry.RetryableError(fmt.Errorf("async rule creation still in progress"))
		case "error":
			errorMsg := "async rule creation failed"
			if status.Error != "" {
				errorMsg = fmt.Sprintf("async rule creation failed: %s", status.Error)
			}
			return retry.NonRetryableError(fmt.Errorf("%s", errorMsg))
		default:
			return retry.NonRetryableError(fmt.Errorf("unexpected async rule creation status: %s", status.Status))
		}
	})
}
