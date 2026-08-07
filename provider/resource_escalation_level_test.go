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

// Regression test for delay being reset to 0 by an update to some other field.
//
// client.EscalationLevel serialises Delay without omitempty so an explicit 0
// reaches the API (#351). That means every update sends a delay, so Update must
// always populate it from state. When it was guarded behind HasChange("delay"),
// changing only the paging strategy left Delay at Go's zero value and silently
// wiped the configured delay.
func TestAccResourceEscalationLevelDelaySurvivesUnrelatedUpdate(t *testing.T) {
	epName := acctest.RandomWithPrefix("tf-ep")
	teamName := acctest.RandomWithPrefix("tf-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationLevelDelayConfig(epName, teamName, "default"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "5"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "paging_strategy_configuration_strategy", "default"),
				),
			},
			{
				// Only the paging strategy changes; delay must be left alone.
				Config: testAccResourceEscalationLevelDelayConfig(epName, teamName, "random"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "5"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "paging_strategy_configuration_strategy", "random"),
				),
			},
		},
	})
}

func testAccResourceEscalationLevelDelayConfig(epName, teamName, strategy string) string {
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
	delay                = 5

	paging_strategy_configuration_strategy = "%s"

	notification_target_params {
		id   = rootly_team.test.id
		type = "team"
	}
}
`, teamName, epName, strategy)
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
