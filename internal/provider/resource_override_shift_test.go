package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
)

func TestAccOverrideShiftResource(t *testing.T) {
	resName := "rootly_override_shift.test"
	scheduleName := acctest.RandomWithPrefix("tf-schedule")

	configStateChecks := []statecheck.StateCheck{
		statecheck.CompareValuePairs(resName, tfjsonpath.New("schedule_id"), "rootly_schedule.test", tfjsonpath.New("id"), compare.ValuesSame()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("user_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("is_override"), knownvalue.Bool(true)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccOverrideShiftResourceConfig(scheduleName),
				ConfigStateChecks: append(configStateChecks),
			},
			{
				Config:            testAccOverrideShiftResourceConfig(scheduleName),
				ConfigStateChecks: append(configStateChecks),
			},
		},
	})
}

func testAccOverrideShiftResourceConfig(scheduleName string) string {
	return testAccUserDataSourceConfig + fmt.Sprintf(`
resource "rootly_schedule" "test" {
	name              = "%s"
	owner_user_id     = data.rootly_user.test.id
	all_time_coverage = true
}

resource "rootly_override_shift" "test" {
	starts_at   = "2025-10-01T15:00:00Z"
	ends_at     = "2025-10-02T03:00:00Z"
	is_override = true
	schedule_id = rootly_schedule.test.id
	user_id     = data.rootly_user.test.id
}
`, scheduleName)
}
