package provider

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestResourceWorkflowTaskHTTPClientUpgradeResourceState(t *testing.T) {
	t.Parallel()

	resourceSchema := resourceWorkflowTaskHttpClient()
	server := schema.NewGRPCProviderServer(&schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"rootly_workflow_task_http_client": resourceSchema,
		},
	})

	testCases := map[string]string{
		"pre-v5.18 string state": `{"workflow_id":"workflow-id","task_params":[{"url":"https://example.com","succeed_on_status":"200","retry_count":"1","retry_wait_time":"15"}]}`,
		"v5.18 numeric state":    `{"workflow_id":"workflow-id","task_params":[{"url":"https://example.com","succeed_on_status":"200","retry_count":1,"retry_wait_time":15}]}`,
	}

	for name, rawState := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
				TypeName: "rootly_workflow_task_http_client",
				Version:  0,
				RawState: &tfprotov5.RawState{JSON: []byte(rawState)},
			})
			require.NoError(t, err)
			require.Empty(t, response.Diagnostics)

			state, err := msgpack.Unmarshal(response.UpgradedState.MsgPack, resourceSchema.CoreConfigSchema().ImpliedType())
			require.NoError(t, err)

			params := state.GetAttr("task_params").Index(cty.NumberIntVal(0))
			require.Equal(t, int64(1), ctyNumberAsInt64(t, params.GetAttr("retry_count")))
			require.Equal(t, int64(15), ctyNumberAsInt64(t, params.GetAttr("retry_wait_time")))
		})
	}
}

func ctyNumberAsInt64(t *testing.T, value cty.Value) int64 {
	t.Helper()

	result, accuracy := value.AsBigFloat().Int64()
	require.Equal(t, big.Exact, accuracy)
	return result
}
