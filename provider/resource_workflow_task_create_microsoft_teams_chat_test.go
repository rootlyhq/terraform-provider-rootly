package provider

// Skipped: API requires real MS Teams members with email addresses and
// at least 2 members for group chats. Cannot test with dummy data.

import "testing"

func TestAccResourceWorkflowTaskCreateMicrosoftTeamsChat(t *testing.T) {
	t.Skip("Requires real MS Teams members with email addresses")
}
