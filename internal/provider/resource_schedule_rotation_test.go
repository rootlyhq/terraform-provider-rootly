package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
)

func TestAccResourceScheduleRotation_Validation(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "rootly_schedule_rotation" "test" {
						schedule_id          = "123"
						name                 = "name"
						active_all_week      = true
						active_time_type     = "all_day"
						position             = 1
						start_time           = "2025-06-20T00:00:00Z"

						schedule_rotationable_type       = "ScheduleCustomRotation"
						schedule_rotationable_attributes = {
							shift_length      = 7
							shift_length_unit = "days"
							handoff_time      = "09:00"
						}

						schedule_rotation_members {
							position    = 1
							member_type = "User"
							member_id   = "1"
						}

						// Duplicate position should fail
						schedule_rotation_members {
							position    = 1
							member_type = "User"
							member_id   = "2"
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralErrors(
					`• Conflicting Path 1: schedule_rotation_members[Value({"member_id":"1","member_type":"User","position":1})].position`,
					`• Conflicting Path 2: schedule_rotation_members[Value({"member_id":"2","member_type":"User","position":1})].position`,
					`Values for "position" must be unique across all elements in the set.`,
				),
			},
		},
	})
}

func TestAccResourceScheduleRotation_UpgradeFromVersion(t *testing.T) {
	addr := "rootly_schedule_rotation.test"
	name := acctest.RandomWithPrefix("tf-rotation")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("position"), knownvalue.Int64Exact(1)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_type"), knownvalue.StringExact("ScheduleCustomRotation")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_all_week"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_type"), knownvalue.StringExact("all_day")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_zone"), knownvalue.StringExact("UTC")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("start_time"), knownvalue.StringExact("2025-06-20T00:00:00Z")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_attributes"), knownvalue.SetSizeExact(0)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotation_members"), knownvalue.SetExact([]knownvalue.Check{
			knownvalue.ObjectExact(map[string]knownvalue.Check{
				"member_id":   knownvalue.NotNull(),
				"member_type": knownvalue.StringExact("User"),
				"position":    knownvalue.Int64Exact(1),
			}),
		})),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rootly": {
						Source:            "rootlyhq/rootly",
						VersionConstraint: "5.17.2",
					},
				},
				Config: testAccResourceScheduleRotationConfig(name, ""),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_attributes"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"handoff_time":      knownvalue.StringExact("09:00"),
						"shift_length":      knownvalue.StringExact("7"),
						"shift_length_unit": knownvalue.StringExact("days"),
					})),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccResourceScheduleRotationConfig(name, ""),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.Null()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_attributes"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"handoff_time":      knownvalue.StringExact("09:00"),
						"handoff_day":       knownvalue.Null(),
						"shift_length":      knownvalue.Int64Exact(7),
						"shift_length_unit": knownvalue.StringExact("days"),
					})),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config: testAccResourceScheduleRotationConfig(name, `
					active_days = ["M", "T", "W"]
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("M"),
						knownvalue.StringExact("T"),
						knownvalue.StringExact("W"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.Null()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_attributes"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"handoff_time":      knownvalue.StringExact("09:00"),
						"handoff_day":       knownvalue.Null(),
						"shift_length":      knownvalue.Int64Exact(7),
						"shift_length_unit": knownvalue.StringExact("days"),
					})),
				),
			},
		},
	})
}

func TestAccResourceScheduleRotation_Basic(t *testing.T) {
	addr := "rootly_schedule_rotation.test"
	name := acctest.RandomWithPrefix("tf-rotation")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("position"), knownvalue.Int64Exact(1)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_type"), knownvalue.StringExact("ScheduleCustomRotation")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_all_week"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_type"), knownvalue.StringExact("all_day")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_zone"), knownvalue.StringExact("UTC")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("start_time"), knownvalue.StringExact("2025-06-20T00:00:00Z")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotationable_attributes"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"handoff_time":      knownvalue.StringExact("09:00"),
			"handoff_day":       knownvalue.Null(),
			"shift_length":      knownvalue.Int64Exact(7),
			"shift_length_unit": knownvalue.StringExact("days"),
		})),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotation_members"), knownvalue.SetExact([]knownvalue.Check{
			knownvalue.ObjectExact(map[string]knownvalue.Check{
				"member_id":   knownvalue.NotNull(),
				"member_type": knownvalue.StringExact("User"),
				"position":    knownvalue.Int64Exact(1),
			}),
		})),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationConfig(name, ""),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_attributes"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.Null()),
				),
			},
			{
				Config: testAccResourceScheduleRotationConfig(name+"-updated", `
					active_days = ["M", "T", "W"]
					end_time = "2025-06-21T00:00:00Z"

					active_time_attributes {
						start_time = "09:00"
						end_time = "12:00"
					}

					active_time_attributes {
						start_time = "18:00"
						end_time = "21:00"
					}

				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("M"),
						knownvalue.StringExact("T"),
						knownvalue.StringExact("W"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_attributes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_time": knownvalue.StringExact("09:00"),
							"end_time":   knownvalue.StringExact("12:00"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_time": knownvalue.StringExact("18:00"),
							"end_time":   knownvalue.StringExact("21:00"),
						}),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.StringExact("2025-06-21T00:00:00Z")),
				),
			},
			{
				Config: testAccResourceScheduleRotationConfig(name+"-updated", `
					active_time_attributes {
						start_time = "13:00"
						end_time = "14:00"
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_days"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("active_time_attributes"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_time": knownvalue.StringExact("13:00"),
							"end_time":   knownvalue.StringExact("14:00"),
						}),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("end_time"), knownvalue.Null()),
				),
			},
		},
	})
}

func testAccResourceScheduleRotationConfig(name string, extras string) string {
	return fmt.Sprintf(`
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "test" {
	name          = "%[1]s"
	owner_user_id = data.rootly_user.test.id
}

resource "rootly_schedule_rotation" "test" {
	schedule_id          = rootly_schedule.test.id
	name                 = "%[1]s"
	active_all_week      = true
	active_time_type     = "all_day"
	position             = 1
	start_time           = "2025-06-20T00:00:00Z"

	schedule_rotationable_type       = "ScheduleCustomRotation"
	schedule_rotationable_attributes = {
		shift_length      = 7
		shift_length_unit = "days"
		handoff_time      = "09:00"
	}

	%[2]s

	time_zone = "UTC"

	schedule_rotation_members {
		position    = 1
		member_type = "User"
		member_id   = data.rootly_user.test.id
	}
}
`, name, extras)
}

func TestAccResourceScheduleRotation_MemberUpdate(t *testing.T) {
	addr := "rootly_schedule_rotation.test"
	name := acctest.RandomWithPrefix("tf-rotation")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleRotationMemberUpdateConfig(name, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotation_members"), knownvalue.SetSizeExact(1)),
				},
			},
			{
				Config: testAccResourceScheduleRotationMemberUpdateConfig(name, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotation_members"), knownvalue.SetSizeExact(2)),
				},
			},
			{
				Config: testAccResourceScheduleRotationMemberUpdateConfig(name, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule_rotation_members"), knownvalue.SetSizeExact(1)),
				},
			},
		},
	})
}

func testAccResourceScheduleRotationMemberUpdateConfig(name string, withSecondMember bool) string {
	secondMember := ""
	if withSecondMember {
		secondMember = `
	schedule_rotation_members {
		position    = 2
		member_type = "Schedule"
		member_id   = rootly_schedule.nested.id
	}
`
	}

	return fmt.Sprintf(`
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}

resource "rootly_schedule" "test" {
	name          = "%[1]s"
	owner_user_id = data.rootly_user.test.id
}

resource "rootly_schedule" "nested" {
	name          = "%[1]s-nested"
	owner_user_id = data.rootly_user.test.id
}

resource "rootly_schedule_rotation" "test" {
	schedule_id          = rootly_schedule.test.id
	name                 = "%[1]s"
	active_all_week      = true
	active_time_type     = "all_day"
	position             = 1
	start_time           = "2025-06-20T00:00:00Z"

	schedule_rotationable_type       = "ScheduleCustomRotation"
	schedule_rotationable_attributes = {
		shift_length      = 7
		shift_length_unit = "days"
		handoff_time      = "09:00"
	}

	time_zone = "UTC"

	schedule_rotation_members {
		position    = 1
		member_type = "User"
		member_id   = data.rootly_user.test.id
	}
%[2]s
}
`, name, secondMember)
}
