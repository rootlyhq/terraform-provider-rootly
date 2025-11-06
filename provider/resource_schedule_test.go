package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/must"
)

func TestAccResourceSchedule(t *testing.T) {
	resName := "rootly_schedule.tf"
	scheduleName := acctest.RandomWithPrefix("tf-schedule")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScheduleConfig(testAccResourceScheduleConfigData{
					Name:            scheduleName,
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
					`,
					RotationName: "test-initial",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", scheduleName),
					resource.TestCheckResourceAttr(resName, "description", "test description"),
					resource.TestCheckResourceAttr(resName, "all_time_coverage", "true"),
					resource.TestCheckResourceAttr(resName, "owner_user_id", "4261"),
					resource.TestCheckResourceAttr(resName, "owner_group_ids.0", "a19ce0d4-8033-410b-97dd-c51164eadfc6"),
					resource.TestCheckResourceAttr(resName, "slack_user_group.id", "123XYZ"),
					resource.TestCheckResourceAttr(resName, "slack_user_group.name", "slack user group"),
					resource.TestCheckResourceAttr(resName, "slack_channel.id", "456ABC"),
					resource.TestCheckResourceAttr(resName, "slack_channel.name", "slack channel"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-initial"),
				),
			},
			{
				Config: testAccResourceScheduleConfig(testAccResourceScheduleConfigData{
					Name:            scheduleName + "-updated",
					Description:     "test updated description",
					OwnerGroupId:    "868f05dd-3c8f-4fe8-8aa7-6c4851b72c15",
					OwnerUserId:     117092,
					AllTimeCoverage: false,
					Extras:          ``,
					RotationName:    "test-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resName, "slack_user_group.id"),
					resource.TestCheckResourceAttr(resName, "name", scheduleName+"-updated"),
					resource.TestCheckResourceAttr(resName, "description", "test updated description"),
					resource.TestCheckResourceAttr(resName, "owner_user_id", "117092"),
					resource.TestCheckResourceAttr(resName, "owner_group_ids.0", "868f05dd-3c8f-4fe8-8aa7-6c4851b72c15"),
					resource.TestCheckResourceAttr(resName, "slack_user_group.#", "0"),
					resource.TestCheckResourceAttr(resName, "slack_channel.#", "0"),
					resource.TestCheckResourceAttr(resName, "all_time_coverage", "false"),
					resource.TestCheckResourceAttr("rootly_schedule_rotation.tf", "name", "test-updated"),
				),
			},
			{
				ResourceName:      resName,
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
