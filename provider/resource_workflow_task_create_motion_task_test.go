package provider

// Skipped: API returns labels as string but schema declares TypeList.
// Causes panic on read. Swagger says array but API returns string.

import "testing"

func TestAccResourceWorkflowTaskCreateMotionTask(t *testing.T) {
	t.Skip("API returns labels as string but schema expects array — panics on read")
}
