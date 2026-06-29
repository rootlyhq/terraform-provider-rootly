package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			heartbeat, err := testAccSharedClient.GetHeartbeatWithResponse(ctx, id)
			if err != nil {
				return fmt.Errorf("failed to get heartbeat: %v", err)
			}

			if heartbeat == nil {
				return fmt.Errorf("heartbeat not found")
			}

			if heartbeat.ApplicationvndApiJSON200.Data.Attributes.Enabled != enabled {
				return fmt.Errorf("heartbeat enabled status does not match expected value: %t, got: %t", enabled, heartbeat.ApplicationvndApiJSON200.Data.Attributes.Enabled)
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
				// Create with an owning team.
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, true, true),
				Check: resource.ComposeTestCheckFunc(
					checkHeartbeatEnabled(true),
					resource.TestCheckResourceAttr("rootly_heartbeat.test", "owner_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_heartbeat.test", "owner_group_ids.0", "rootly_team.test", "id"),
				),
			},
			{
				// Update: clear the owning team.
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, false, false),
				Check: resource.ComposeTestCheckFunc(
					checkHeartbeatEnabled(false),
					resource.TestCheckResourceAttr("rootly_heartbeat.test", "owner_group_ids.#", "0"),
				),
			},
			{
				// Update: re-add the owning team.
				Config: testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName, true, true),
				Check: resource.ComposeTestCheckFunc(
					checkHeartbeatEnabled(true),
					resource.TestCheckResourceAttr("rootly_heartbeat.test", "owner_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_heartbeat.test", "owner_group_ids.0", "rootly_team.test", "id"),
				),
			},
		},
	})
}

func testAccResourceHeartbeat(heartbeatName, teamName, alertUrgencyName string, heartbeatEnabled, withOwnerTeam bool) string {
	ownerGroupIds := "[]"
	if withOwnerTeam {
		ownerGroupIds = "[rootly_team.test.id]"
	}

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
  owner_group_ids = %s
  enabled = %t
}
`, teamName, alertUrgencyName, heartbeatName, ownerGroupIds, heartbeatEnabled)
}
