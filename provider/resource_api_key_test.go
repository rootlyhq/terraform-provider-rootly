package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceApiKey(t *testing.T) {
	t.Skip("Skipping: CI uses a service account token which cannot create personal API keys")
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyConfig(rName, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
				),
			},
		},
	})
}

func TestAccResourceApiKeyTeamScoped(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)
	var firstID string

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyTeamScopedConfig(rName, expiresAt, "first"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_api_key.test", "kind", "team"),
					resource.TestCheckResourceAttrPair("rootly_api_key.test", "group_id", "rootly_team.first", "id"),
					testAccCaptureResourceID("rootly_api_key.test", &firstID),
				),
			},
			{
				ResourceName:      "rootly_api_key.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceApiKeyTeamScopedConfig(rName, expiresAt, "second"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("rootly_api_key.test", "group_id", "rootly_team.second", "id"),
					testAccCheckResourceIDChanged("rootly_api_key.test", &firstID),
				),
			},
		},
	})
}

func testAccCaptureResourceID(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckResourceIDChanged(name string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		if rs.Primary.ID == *previous {
			return fmt.Errorf("expected %s to be replaced, but ID %s is unchanged", name, rs.Primary.ID)
		}
		return nil
	}
}

func testAccResourceApiKeyConfig(name, expiresAt string) string {
	return fmt.Sprintf(`
resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "%s"
}
`, name, expiresAt)
}

func testAccResourceApiKeyTeamScopedConfig(name, expiresAt, team string) string {
	return fmt.Sprintf(`
resource "rootly_team" "first" {
	name = "%[1]s-first"
}

resource "rootly_team" "second" {
	name = "%[1]s-second"
}

resource "rootly_api_key" "test" {
	name       = "%[1]s"
	expires_at = "%[2]s"
	kind       = "team"
	group_id   = rootly_team.%[3]s.id
}
`, name, expiresAt, team)
}
