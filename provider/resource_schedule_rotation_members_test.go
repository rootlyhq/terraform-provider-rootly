package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScheduleRotationMembers(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-sched-rm")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationMembersCreatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", rName+"-rotation"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "User"),
					resource.TestCheckResourceAttrSet("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
				),
			},
			{
				Config: testAccResourceScheduleRotationMembersUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", rName+"-rotation-updated"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "User"),
					resource.TestCheckResourceAttrSet("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
				),
			},
		},
	})
}

func testAccResourceScheduleRotationMembersCreatedConfig(rName string) string {
	return fmt.Sprintf(`
data "rootly_user" "test_user" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "tf" {
	name = "%s-schedule"
	description = "test schedule for rotation members"
	owner_user_id = data.rootly_user.test_user.id
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "%s-rotation"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 7
		shift_length_unit = "days"
		handoff_time = "09:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "UTC"

	schedule_rotation_members {
		member_type = "User"
		member_id   = data.rootly_user.test_user.id
		position    = 1
	}
}
`, rName, rName)
}

func testAccResourceScheduleRotationMembersUpdatedConfig(rName string) string {
	return fmt.Sprintf(`
data "rootly_user" "test_user" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "tf" {
	name = "%s-schedule"
	description = "test schedule for rotation members updated"
	owner_user_id = data.rootly_user.test_user.id
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "%s-rotation-updated"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 7
		shift_length_unit = "days"
		handoff_time = "09:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "UTC"

	schedule_rotation_members {
		member_type = "User"
		member_id   = data.rootly_user.test_user.id
		position    = 1
	}
}
`, rName, rName)
}

func TestAccResourceScheduleRotationMembersScheduleType(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-sched-rm-s")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationMembersWithScheduleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", rName+"-rotation"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "Schedule"),
					resource.TestCheckResourceAttrSet("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
				),
			},
		},
	})
}

func testAccResourceScheduleRotationMembersWithScheduleConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_schedule" "parent_tf" {
	name = "%s-parent"
	description = "parent schedule for nesting test"
	all_time_coverage = true
}

resource "rootly_schedule" "nested_tf" {
	name = "%s-nested"
	description = "nested schedule for testing"
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.parent_tf.id
	name            = "%s-rotation"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 7
		shift_length_unit = "days"
		handoff_time = "09:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "UTC"

	schedule_rotation_members {
		member_type = "Schedule"
		member_id   = rootly_schedule.nested_tf.id
		position    = 1
	}
}
`, rName, rName, rName)
}

func TestAccResourceScheduleRotationMembersValidation(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-sched-rm-v")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScheduleRotationMembersInvalidTypeConfig(rName),
				ExpectError: regexp.MustCompile(`expected.*member_type.*to be one of.*got InvalidType`),
			},
		},
	})
}

func testAccResourceScheduleRotationMembersInvalidTypeConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_schedule" "tf" {
	name = "%s-schedule"
	description = "test schedule for invalid member type"
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "%s-rotation"
	active_all_week = true
	active_time_type = "all_day"
	position         = 1
	start_time       = "2025-06-20T00:00:00Z"
	schedule_rotationable_attributes = {
		shift_length = 7
		shift_length_unit = "days"
		handoff_time = "09:00"
	}
	schedule_rotationable_type = "ScheduleCustomRotation"
	time_zone                  = "UTC"

	schedule_rotation_members {
		member_type = "InvalidType"
		member_id   = "dummy-id"
		position    = 1
	}
}
`, rName, rName)
}
