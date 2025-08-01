package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSchedule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleCreated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule.tf", "name", "test-initial"),
					resource.TestCheckResourceAttr("rootly_schedule.tf", "slack_user_group.id", "123XYZ"),
					resource.TestCheckResourceAttr("rootly_schedule.tf", "slack_user_group.name", "test"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-initial"),
				),
			},
			{
				Config: testAccResourceScheduleUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("rootly_schedule.tf", "slack_user_group.id"),
					resource.TestCheckResourceAttr("rootly_schedule.tf", "name", "test-updated"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-updated"),
				),
			},
		},
	})
}

const testAccResourceScheduleCreated = `
resource "rootly_schedule" "tf" {
	name = "test-initial"
	slack_user_group = {
		id = "123XYZ"
		name = "test"
	}
}
resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "test-initial"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 5
		shift_length_unit = "days"
		handoff_time = "12:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "Europe/Athens"
}
`

const testAccResourceScheduleUpdated = `
resource "rootly_schedule" "tf" {
	name = "test-updated"
}
resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "test-updated"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-07-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 5
		shift_length_unit = "days"
		handoff_time = "12:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "Europe/Athens"
}
`
