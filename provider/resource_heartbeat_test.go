package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rootly "github.com/rootlyhq/rootly-go"
)

func TestAccResourceHeartbeat(t *testing.T) {
	heartbeatName := acctest.RandomWithPrefix("tf-heartbeat")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")

	checkHeartbeatEnabled := func(enabled bool) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources["rootly_heartbeat.test"]
			if !ok {
				return fmt.Errorf("not found: rootly_heartbeat.test")
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("ID is not set")
			}

			id := rs.Primary.ID
			ctx := context.Background()
			heartbeat, err := testAccSharedClient.GetHeartbeatWithResponse(ctx, rootly.ID(id))
			if err != nil {
				return fmt.Errorf("failed to get heartbeat: %v", err)
			}

			if heartbeat == nil {
				return fmt.Errorf("heartbeat not found")
			}

			if heartbeat.ApplicationVndAPIJSON200.Data.Attributes.Enabled != enabled {
				return fmt.Errorf("heartbeat enabled status does not match expected value: %t, got: %t", enabled, heartbeat.ApplicationVndAPIJSON200.Data.Attributes.Enabled)
			}

			return nil
		}
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, true),
				Check:  checkHeartbeatEnabled(true),
			},
			{
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, false),
				Check:  checkHeartbeatEnabled(false),
			},
			{
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, true),
				Check:  checkHeartbeatEnabled(true),
			},
		},
	})
}

func testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName string, heartbeatEnabled bool) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
  name = "%s"
}

resource "rootly_alert_urgency" "test" {
  name = "%s"
  description = "Test alert urgency"
  position = 1
}

resource "rootly_heartbeat" "test" {
  name = "%s"
  alert_summary = "Heartbeat expired"
  interval = 5
  notification_target_id = rootly_team.test.id
  notification_target_type = "Group"
  alert_urgency_id = rootly_alert_urgency.test.id
  enabled = %t
}
`, teamName, alertUrgencyName, heartbeatName, heartbeatEnabled)
}
