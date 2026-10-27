package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialAttributesAreSensitive(t *testing.T) {
	t.Parallel()

	heartbeat := resourceHeartbeat().Schema
	require.True(t, heartbeat["secret"].Sensitive, "rootly_heartbeat.secret is a bearer token")
	require.True(t, heartbeat["ping_url"].Sensitive, "rootly_heartbeat.ping_url accepts unauthenticated pings")

	statusPage := resourceStatusPage().Schema
	require.True(t, statusPage["authentication_password"].Sensitive, "rootly_status_page.authentication_password is a password")
}
