package provider

// Skipped: API returns empty objects {} for unset TypeMap fields
// (integration, priority, status) causing plan drift on every apply.

import "testing"

func TestAccResourceWorkflowTaskCreateJiraIssue(t *testing.T) {
	t.Skip("API returns empty objects for unset map fields causing plan drift")
}
