// Package polling provides utilities for polling operations until they complete
package polling

import (
	"github.com/google/uuid"
)

// GenerateRequestId creates a new UUID for tracking async operations
func GenerateRequestId() string {
	return uuid.New().String()
}
