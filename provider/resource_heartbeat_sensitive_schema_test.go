package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceHeartbeatSensitiveSchema(t *testing.T) {
	t.Parallel()

	resourceSchema := resourceHeartbeat().Schema

	require.True(t, resourceSchema["secret"].Sensitive, "secret is a bearer token")
	require.True(t, resourceSchema["ping_url"].Sensitive, "ping_url accepts pings without the secret")
}
