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

// Guards both halves of the delay contract end to end:
//
//   - step 2 changes only notification_target_params. delay is not in the diff, so the update
//     payload must omit it and the server must keep 15. Before Delay became a nullable *int
//     the zero value was serialized here and reset delay to 0.
//   - step 3 changes delay to 0 explicitly. That must actually be applied, not swallowed by
//     omitempty (TER-182, #351).
func TestAccResourceEscalationLevelDelayPartialUpdate(t *testing.T) {
	epName := acctest.RandomWithPrefix("tf-ep")
	teamName := acctest.RandomWithPrefix("tf-team")
	team2Name := acctest.RandomWithPrefix("tf-team")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationLevelDelayConfig(epName, teamName, team2Name, 15, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "15"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "notification_target_params.#", "1"),
				),
			},
			{
				Config: testAccResourceEscalationLevelDelayConfig(epName, teamName, team2Name, 15, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "15"),
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "notification_target_params.#", "2"),
				),
			},
			{
				Config: testAccResourceEscalationLevelDelayConfig(epName, teamName, team2Name, 0, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_escalation_level.test", "delay", "0"),
				),
			},
		},
	})
}

func testAccResourceEscalationLevelDelayConfig(epName, teamName, team2Name string, delay int, secondTarget bool) string {
	extraTarget := ""
	if secondTarget {
		extraTarget = `
	notification_target_params {
		id   = rootly_team.test2.id
		type = "team"
	}
`
	}

	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}

resource "rootly_team" "test2" {
	name = "%s"
}

resource "rootly_escalation_policy" "test" {
	name = "%s"
}

resource "rootly_escalation_level" "test" {
	escalation_policy_id = rootly_escalation_policy.test.id
	position             = 1
	delay                = %d

	notification_target_params {
		id   = rootly_team.test.id
		type = "team"
	}
%s
}
`, teamName, team2Name, epName, delay, extraTarget)
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
