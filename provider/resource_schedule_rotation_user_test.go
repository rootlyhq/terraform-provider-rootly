package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScheduleRotationUser(t *testing.T) {
	scheduleName := acctest.RandomWithPrefix("tf-sched")
	rotationName := acctest.RandomWithPrefix("tf-rotation")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationUserConfig(scheduleName, rotationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_schedule_rotation_user.test", "id"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation_user.test", "position", "1"),
				),
			},
		},
	})
}

func testAccResourceScheduleRotationUserConfig(scheduleName, rotationName string) string {
	return fmt.Sprintf(`
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "test" {
	name = "%s"
}

resource "rootly_schedule_rotation" "test" {
	schedule_id     = rootly_schedule.test.id
	name            = "%s"
	active_all_week  = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length      = 7
		shift_length_unit = "days"
		handoff_time      = "09:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "UTC"
}

resource "rootly_schedule_rotation_user" "test" {
	schedule_rotation_id = rootly_schedule_rotation.test.id
	user_id              = data.rootly_user.test.id
	position             = 1
}
`, scheduleName, rotationName)
}
