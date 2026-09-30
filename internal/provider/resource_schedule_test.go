package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/must"
)

func TestAccResourceSchedule_Basic(t *testing.T) {
	addr := "rootly_schedule.tf"
	name := acctest.RandomWithPrefix("tf-schedule")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleConfig(testAccResourceScheduleConfigData{
					Name:            name,
					Description:     "test description",
					OwnerGroupId:    "a19ce0d4-8033-410b-97dd-c51164eadfc6",
					OwnerUserId:     4261,
					AllTimeCoverage: true,
					Extras: `
						slack_user_group = {
							id = "123XYZ"
							name = "slack user group"
						}
						slack_channel = {
							id = "456ABC"
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
					`,
					RotationName: "test-initial",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("description"), knownvalue.StringExact("test description")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("all_time_coverage"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("owner_user_id"), knownvalue.Int64Exact(4261)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("owner_group_ids"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.StringExact("a19ce0d4-8033-410b-97dd-c51164eadfc6"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_user_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.StringExact("123XYZ"),
						"name": knownvalue.StringExact("slack user group"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_channel"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.StringExact("456ABC"),
						"name": knownvalue.StringExact("slack channel"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("sync_linear_enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("include_shadows_in_slack_notifications"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_start_notifications_enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_update_notifications_enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_day_of_week"), knownvalue.StringExact("tuesday")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_time_of_day"), knownvalue.StringExact("10:30")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_time_zone"), knownvalue.StringExact("Australia/Sydney")),
				},
			},
			{
				Config: testAccResourceScheduleConfig(testAccResourceScheduleConfigData{
					Name:            name + "-updated",
					Description:     "test updated description",
					OwnerGroupId:    "868f05dd-3c8f-4fe8-8aa7-6c4851b72c15",
					OwnerUserId:     117092,
					AllTimeCoverage: false,
					Extras: `
						slack_user_group = {
							id = "789DEF"
							name = "updated slack user group"
						}
						slack_channel = {
							id = "012GHI"
							name = "updated slack channel"
						}
						sync_linear_enabled                    = false
						include_shadows_in_slack_notifications = false
						shift_start_notifications_enabled      = false
						shift_update_notifications_enabled     = false
						shift_report_enabled                   = false
					`,
					RotationName: "test-updated",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("description"), knownvalue.StringExact("test updated description")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("all_time_coverage"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("owner_user_id"), knownvalue.Int64Exact(117092)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("owner_group_ids"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.StringExact("868f05dd-3c8f-4fe8-8aa7-6c4851b72c15"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_user_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.StringExact("789DEF"),
						"name": knownvalue.StringExact("updated slack user group"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_channel"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.StringExact("012GHI"),
						"name": knownvalue.StringExact("updated slack channel"),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("sync_linear_enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("include_shadows_in_slack_notifications"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_start_notifications_enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_update_notifications_enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_day_of_week"), knownvalue.StringExact("tuesday")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_time_of_day"), knownvalue.StringExact("10:30")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("shift_report_time_zone"), knownvalue.StringExact("Australia/Sydney")),
				},
			},
			{
				Config: testAccResourceScheduleConfig(testAccResourceScheduleConfigData{
					Name:            name + "-null-test",
					Description:     "test null values",
					OwnerGroupId:    "868f05dd-3c8f-4fe8-8aa7-6c4851b72c15",
					OwnerUserId:     117092,
					AllTimeCoverage: false,
					Extras: `
						slack_user_group = {}
						slack_channel = {}
						sync_linear_enabled = false
					`,
					RotationName: "test-null",
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-null-test")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("description"), knownvalue.StringExact("test null values")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_user_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.Null(),
						"name": knownvalue.Null(),
					})),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("slack_channel"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.Null(),
						"name": knownvalue.Null(),
					})),
				},
			},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccResourceScheduleConfigTemplate = template.Must(template.New("config").Parse(`
resource "rootly_schedule" "tf" {
	name = "{{ .Name }}"
	description = "{{ .Description }}"
	owner_group_ids = ["{{ .OwnerGroupId }}"]
	owner_user_id = {{ .OwnerUserId }}
	all_time_coverage = {{ .AllTimeCoverage }}
	{{ .Extras }}
}

resource "rootly_schedule_rotation" "tf" {
	schedule_id     = rootly_schedule.tf.id
	name            = "{{ .RotationName }}"
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
`))

type testAccResourceScheduleConfigData struct {
	Name            string
	Description     string
	OwnerGroupId    string
	OwnerUserId     int
	AllTimeCoverage bool
	Extras          string
	RotationName    string
}

func testAccResourceScheduleConfig(data testAccResourceScheduleConfigData) string {
	var buf bytes.Buffer
	must.Do(testAccResourceScheduleConfigTemplate.Execute(&buf, data))
	return buf.String()
}
