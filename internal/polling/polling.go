// Package polling provides utilities for polling operations until they complete
package polling

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// AsyncRuleCreationStatus represents the response from async rule creation status endpoint
type AsyncRuleCreationStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// AsyncRuleCreationStatusFunc defines the function signature for checking async rule creation status
type AsyncRuleCreationStatusFunc func(alertRouteID, jobID string) (*AsyncRuleCreationStatus, error)

// RefetchAlertRouteFunc defines the function signature for refetching an alert route
type RefetchAlertRouteFunc func() error

// GenerateJobID creates a new UUID for tracking async operations
func GenerateJobID() string {
	return uuid.New().String()
}

// WaitForAsyncRuleCreation polls the async rule creation status endpoint until completion
func WaitForAsyncRuleCreation(ctx context.Context, alertRouteID, jobID string, 
	checkStatusFunc AsyncRuleCreationStatusFunc, refetchFunc RefetchAlertRouteFunc) error {
	
	return retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		status, err := checkStatusFunc(alertRouteID, jobID)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error checking async rule creation status: %w", err))
		}

		switch status.Status {
		case "success":
			// Rules creation completed successfully, now refetch the alert route
			if err := refetchFunc(); err != nil {
				return retry.NonRetryableError(fmt.Errorf("error refetching alert route after successful rule creation: %w", err))
			}
			return nil // Success
		case "pending":
			return retry.RetryableError(fmt.Errorf("async rule creation still in progress"))
		case "error":
			errorMsg := "async rule creation failed"
			if status.Error != "" {
				errorMsg = fmt.Sprintf("async rule creation failed: %s", status.Error)
			}
			return retry.NonRetryableError(fmt.Errorf(errorMsg))
		default:
			return retry.NonRetryableError(fmt.Errorf("unexpected async rule creation status: %s", status.Status))
		}
	})
}

