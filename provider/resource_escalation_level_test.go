package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEscalationLevel(t *testing.T) {
	epName := acctest.RandomWithPrefix("tf-ep")
	teamName := acctest.RandomWithPrefix("tf-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationLevelConfig(epName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "position", "1"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "0"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "notification_target_params.#", "1"),
				),
			},
		},
	})
}

func testAccResourceEscalationLevelConfig(epName, teamName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}

resource "rootly_escalation_policy" "test" {
	name = "%s"
}

resource "rootly_escalation_level" "test" {
	escalation_policy_id = rootly_escalation_policy.test.id
	position             = 1
	delay                = 0

	notification_target_params {
		id   = rootly_team.test.id
		type = "team"
	}
}
`, teamName, epName)
}

func TestAccResourceEscalationLevelCycleRoundRobin(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-el")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationLevelCycleConfig(rName, "users", 3, "active_rotation"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_strategy", "cycle"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_schedule_strategy", "everyone"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_repeats_mode", "users"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_repeats", "3"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_rotation_scope", "active_rotation"),
				),
			},
			{
				Config: testAccResourceEscalationLevelCycleConfig(rName, "all", 1, "entire_schedule"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_repeats_mode", "all"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_repeats", "1"),
					resource.TestCheckResourceAttr("rootly_escalation_level.cycle",
						"paging_strategy_configuration_rotation_scope", "entire_schedule"),
				),
			},
		},
	})
}

func testAccResourceEscalationLevelCycleConfig(rName string, repeatsMode string, repeats int, rotationScope string) string {
	return fmt.Sprintf(`
resource "rootly_team" "cycle_test" {
	name = "%s-team"
}

resource "rootly_escalation_policy" "cycle_test" {
	name = "%s-ep"
}

resource "rootly_escalation_level" "cycle" {
	escalation_policy_id                            = rootly_escalation_policy.cycle_test.id
	position                                        = 1
	delay                                           = 5
	paging_strategy_configuration_strategy          = "cycle"
	paging_strategy_configuration_schedule_strategy = "everyone"
	paging_strategy_configuration_repeats_mode      = "%s"
	paging_strategy_configuration_repeats           = %d
	paging_strategy_configuration_rotation_scope    = "%s"

	notification_target_params {
		id   = rootly_team.cycle_test.id
		type = "team"
	}
}
`, rName, rName, repeatsMode, repeats, rotationScope)
}
