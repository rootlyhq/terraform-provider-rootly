package stateupgrade

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgradeWorkflowTaskHTTPClientV0ToV1(t *testing.T) {
	t.Parallel()

	rawState := map[string]any{
		"task_params": []any{
			map[string]any{
				"retry_count":     "1",
				"retry_wait_time": "15",
			},
		},
	}

	upgraded, err := UpgradeWorkflowTaskHTTPClientV0ToV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	params := upgraded["task_params"].([]any)[0].(map[string]any)
	require.Equal(t, 1, params["retry_count"])
	require.Equal(t, 15, params["retry_wait_time"])
}

func TestUpgradeWorkflowTaskHTTPClientV0ToV1DropsEmptyValues(t *testing.T) {
	t.Parallel()

	rawState := map[string]any{
		"task_params": []any{
			map[string]any{
				"retry_count":     "",
				"retry_wait_time": "",
			},
		},
	}

	upgraded, err := UpgradeWorkflowTaskHTTPClientV0ToV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	params := upgraded["task_params"].([]any)[0].(map[string]any)
	require.NotContains(t, params, "retry_count")
	require.NotContains(t, params, "retry_wait_time")
}

func TestUpgradeWorkflowTaskHTTPClientV0ToV1DropsInvalidValues(t *testing.T) {
	t.Parallel()

	rawState := map[string]any{
		"task_params": []any{
			map[string]any{"retry_count": "one"},
		},
	}

	upgraded, err := UpgradeWorkflowTaskHTTPClientV0ToV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	params := upgraded["task_params"].([]any)[0].(map[string]any)
	require.NotContains(t, params, "retry_count")
}
