package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScheduleRotationMembers(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationMembersCreated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-rotation-with-members"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "User"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id", "4261"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.1.member_type", "User"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.1.member_id", "117092"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.1.position", "2"),
				),
			},
			{
				Config: testAccResourceScheduleRotationMembersUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-rotation-updated"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "User"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id", "117092"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
				),
			},
		},
	})
}

func TestAccResourceScheduleRotationMembersScheduleType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationMembersWithSchedule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-rotation-with-schedule"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_type", "Schedule"),
					resource.TestCheckResourceAttrSet("rootly_schedule_rotation.tf", "schedule_rotation_members.0.member_id"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "schedule_rotation_members.0.position", "1"),
				),
			},
		},
	})
}

func TestAccResourceScheduleRotationMembersValidation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScheduleRotationMembersInvalidType,
				ExpectError: regexp.MustCompile(`expected.*member_type.*to be one of.*got InvalidType`),
			},
		},
	})
}

const testAccResourceScheduleRotationMembersCreated = `
resource "rootly_schedule" "tf" {
	name = "test-schedule-for-rotation-members"
	description = "test schedule for rotation members"
	owner_user_id = 4261
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "test-rotation-with-members"
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
		member_id   = "4261"
		position    = 1
	}

	schedule_rotation_members {
		member_type = "User"
		member_id   = "117092"
		position    = 2
	}
}
`

const testAccResourceScheduleRotationMembersUpdated = `
resource "rootly_schedule" "tf" {
	name = "test-schedule-for-rotation-members"
	description = "test schedule for rotation members updated"
	owner_user_id = 117092
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "test-rotation-updated"
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
		member_id   = "117092"
		position    = 1
	}
}
`

const testAccResourceScheduleRotationMembersWithSchedule = `
resource "rootly_schedule" "parent_tf" {
	name = "parent-schedule-for-nested"
	description = "parent schedule for nesting test"
	owner_user_id = 4261
	all_time_coverage = true
}

resource "rootly_schedule" "nested_tf" {
	name = "nested-schedule"
	description = "nested schedule for testing"
	owner_user_id = 4261
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.parent_tf.id
	name            = "test-rotation-with-schedule"
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
`

const testAccResourceScheduleRotationMembersInvalidType = `
resource "rootly_schedule" "tf" {
	name = "test-schedule-for-invalid-type"
	description = "test schedule for invalid member type"
	owner_user_id = 4261
	all_time_coverage = true
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "test-rotation-invalid"
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
		member_id   = "4261"
		position    = 1
	}
}
`