package provider

import "testing"

func TestAccResourceWorkflowTaskUpdateSlackCanvas(t *testing.T) {
	testAccWorkflowTaskSlackCanvas(t, "update_slack_canvas")
}
