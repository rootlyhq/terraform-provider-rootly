package provider

import "testing"

func TestAccResourceWorkflowTaskCreateSlackCanvas(t *testing.T) {
	testAccWorkflowTaskSlackCanvas(t, "create_slack_canvas")
}
