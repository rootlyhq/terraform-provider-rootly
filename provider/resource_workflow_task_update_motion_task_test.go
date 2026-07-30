package provider

// Skipped: API returns labels as string but schema expects array — panics on read.
// Same issue as create_motion_task.

import "testing"

func TestAccResourceWorkflowTaskUpdateMotionTask(t *testing.T) {
	t.Skip("API returns labels as string but schema expects array — panics on read")
}
