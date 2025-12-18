package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/must"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
)

func TestAccResourceOverrideShift(t *testing.T) {
	resName := "rootly_override_shift.test"
	scheduleName := acctest.RandomWithPrefix("tf-schedule")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOverrideShiftConfig(scheduleName, "2025-10-01T09:00:00Z", "2025-10-01T17:00:00Z"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("schedule_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("starts_at"), knownvalue.StringFunc(func(value string) error {
						if !must.Get(time.Parse(time.RFC3339, value)).Equal(must.Get(time.Parse(time.RFC3339, "2025-10-01T09:00:00Z"))) {
							return fmt.Errorf("starts_at is not equal to 2025-10-01T09:00:00Z")
						}
						return nil
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("ends_at"), knownvalue.StringFunc(func(value string) error {
						if !must.Get(time.Parse(time.RFC3339, value)).Equal(must.Get(time.Parse(time.RFC3339, "2025-10-01T17:00:00Z"))) {
							return fmt.Errorf("ends_at is not equal to 2025-10-01T17:00:00Z")
						}
						return nil
					})),
				),
			},
			{
				Config: testAccResourceOverrideShiftConfig(scheduleName, "2025-10-01T10:00:00Z", "2025-10-01T18:00:00Z"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("schedule_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("starts_at"), knownvalue.StringFunc(func(value string) error {
						if !must.Get(time.Parse(time.RFC3339, value)).Equal(must.Get(time.Parse(time.RFC3339, "2025-10-01T10:00:00Z"))) {
							return fmt.Errorf("starts_at is not equal to 2025-10-01T10:00:00Z")
						}
						return nil
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("ends_at"), knownvalue.StringFunc(func(value string) error {
						if !must.Get(time.Parse(time.RFC3339, value)).Equal(must.Get(time.Parse(time.RFC3339, "2025-10-01T18:00:00Z"))) {
							return fmt.Errorf("ends_at is not equal to 2025-10-01T18:00:00Z")
						}
						return nil
					})),
				),
			},
		},
	})
}

func testAccResourceOverrideShiftConfig(scheduleName, startsAt, endsAt string) string {
	return testAccDataSourceUserConfig + fmt.Sprintf(`
resource "rootly_schedule" "test" {
  name          = "%s"
  owner_user_id = data.rootly_user.test.id
}

resource "rootly_override_shift" "test" {
  starts_at   = "%s"
  ends_at     = "%s"
  schedule_id = rootly_schedule.test.id
  user_id     = data.rootly_user.test.id
}
`, scheduleName, startsAt, endsAt)
}
