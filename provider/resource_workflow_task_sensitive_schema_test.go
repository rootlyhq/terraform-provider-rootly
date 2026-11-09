package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestResourceWorkflowTaskSensitiveSchema(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		resource *schema.Resource
		field    string
	}{
		{"create_zoom_meeting", resourceWorkflowTaskCreateZoomMeeting(), "password"},
		{"create_webex_meeting", resourceWorkflowTaskCreateWebexMeeting(), "password"},
		{"http_client", resourceWorkflowTaskHttpClient(), "headers"},
	}

	for _, tc := range cases {
		taskParams := tc.resource.Schema["task_params"].Elem.(*schema.Resource).Schema
		require.True(t, taskParams[tc.field].Sensitive, "%s task_params.%s carries credentials", tc.name, tc.field)
	}
}
