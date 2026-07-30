package provider

// Skipped: API returns 500 on task creation (times out after 9 retries).

import "testing"

func TestAccResourceWorkflowTaskUpdateConfluencePage(t *testing.T) {
	t.Skip("API returns 500 on workflow task creation")
}
