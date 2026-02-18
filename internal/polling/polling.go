// Package polling provides utilities for polling operations until they complete
package polling

import (
	"github.com/google/uuid"
)

func GenerateRequestId() string {
	return uuid.New().String()
}
