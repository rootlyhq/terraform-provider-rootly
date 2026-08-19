package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestResourceEscalationPolicyAcceptsLegacyBusinessHoursDaysList(t *testing.T) {
	t.Parallel()

	server := schema.NewGRPCProviderServer(&schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"rootly_escalation_policy": resourceEscalationPolicy(),
		},
	})

	response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
		TypeName: "rootly_escalation_policy",
		Version:  0,
		RawState: &tfprotov5.RawState{JSON: []byte(`{
			"name":"policy",
			"business_hours":[{"days":["M","T"]}]
		}`)},
	})
	require.NoError(t, err)
	require.Empty(t, response.Diagnostics)
}
