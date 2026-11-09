package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceStatusPageSensitiveSchema(t *testing.T) {
	t.Parallel()

	resourceSchema := resourceStatusPage().Schema

	require.True(t, resourceSchema["authentication_password"].Sensitive, "authentication_password is a password")
}
