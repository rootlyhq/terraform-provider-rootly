package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestResourceWebhooksEndpointSensitiveSchema(t *testing.T) {
	t.Parallel()

	resourceSchema := resourceWebhooksEndpoint().Schema
	customHeader := resourceSchema["custom_headers"].Elem.(*schema.Resource).Schema

	require.True(t, resourceSchema["secret"].Sensitive, "secret signs webhook deliveries")
	require.True(t, customHeader["value"].Sensitive, "custom header values often carry auth tokens")
}
