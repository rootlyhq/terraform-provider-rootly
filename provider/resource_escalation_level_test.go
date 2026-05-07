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

	notification_target_params {
		id   = rootly_team.test.id
		type = "team"
	}
}
`, teamName, epName)
}
