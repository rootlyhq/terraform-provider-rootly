package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceApiKey(t *testing.T) {
	if os.Getenv("ROOTLY_TEST_API_KEY") == "" {
		t.Skip("Skipping: CI uses a service account token which cannot create personal API keys")
	}
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
				ResourceName:            "rootly_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
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

func TestAccResourceApiKeyRotation(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)
	var firstID, firstToken string

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyRotationConfig(rName, expiresAt, "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_api_key.test", "token"),
					resource.TestCheckResourceAttr("rootly_api_key.test", "rotation_trigger", "1"),
					testAccCaptureResourceID("rootly_api_key.test", &firstID),
					testAccCaptureResourceAttr("rootly_api_key.test", "token", &firstToken),
				),
			},
			{
				ResourceName:            "rootly_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "rotation_trigger"},
			},
			{
				Config: testAccResourceApiKeyRotationConfig(rName, expiresAt, "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "rotation_trigger", "2"),
					testAccCheckResourceIDChanged("rootly_api_key.test", &firstID),
					testAccCheckResourceAttrChanged("rootly_api_key.test", "token", &firstToken),
				),
			},
		},
	})
}

func testAccCaptureResourceAttr(name, key string, value *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		*value = rs.Primary.Attributes[key]
		return nil
	}
}

func testAccCheckResourceAttrChanged(name, key string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		if rs.Primary.Attributes[key] == *previous {
			return fmt.Errorf("expected %s.%s to change, but it is unchanged", name, key)
		}
		return nil
	}
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

func testAccResourceApiKeyRotationConfig(name, expiresAt, version string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%[1]s"
}

resource "rootly_api_key" "test" {
	name       = "%[1]s"
	expires_at = "%[2]s"
	kind       = "team"
	group_id   = rootly_team.test.id

	rotation_trigger = "%[3]s"
}
`, name, expiresAt, version)
}
