package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceSchedule looks up a schedule by name and verifies that the
// notifications-tab attributes round-trip from resource → API → data source.
func TestAccDataSourceSchedule(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-ds-schedule")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScheduleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule.test", "name", name),
					resource.TestCheckResourceAttrPair(
						"data.rootly_schedule.by_name", "id",
						"rootly_schedule.test", "id",
					),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "sync_linear_enabled", "true"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "include_shadows_in_slack_notifications", "true"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_start_notifications_enabled", "true"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_update_notifications_enabled", "true"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_report_enabled", "true"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_report_day_of_week", "tuesday"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_report_time_of_day", "10:30"),
					resource.TestCheckResourceAttr("data.rootly_schedule.by_name", "shift_report_time_zone", "Australia/Sydney"),
				),
			},
		},
	})
}

func testAccDataSourceScheduleConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_schedule" "test" {
	name              = "%s"
	owner_user_id     = 4261
	all_time_coverage = true

	slack_channel = {
		id   = "456ABC"
		name = "slack channel"
	}

	sync_linear_enabled                    = true
	include_shadows_in_slack_notifications = true
	shift_start_notifications_enabled      = true
	shift_update_notifications_enabled     = true
	shift_report_enabled                   = true
	shift_report_day_of_week               = "tuesday"
	shift_report_time_of_day               = "10:30"
	shift_report_time_zone                 = "Australia/Sydney"
}

resource "rootly_schedule_rotation" "test" {
	schedule_id      = rootly_schedule.test.id
	name             = "%s-rotation"
	active_all_week  = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length      = 5
		shift_length_unit = "days"
		handoff_time      = "12:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "Europe/Athens"
}

data "rootly_schedule" "by_name" {
	name       = rootly_schedule.test.name
	depends_on = [rootly_schedule.test]
}
`, name, name)
}
