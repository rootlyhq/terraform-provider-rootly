package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceSecretSensitiveSchema(t *testing.T) {
	t.Parallel()

	require.True(t, resourceSecret().Schema["secret"].Sensitive, "secret holds the credential value")
}
