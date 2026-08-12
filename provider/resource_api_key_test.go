package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func testAccResourceApiKeyConfig(name, expiresAt string) string {
	return fmt.Sprintf(`
resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "%s"
}
`, name, expiresAt)
}

func TestAccResourceApiKeyTeamScoped(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-apikey")
	expiresAt := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiKeyTeamScopedConfig(rName, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_api_key.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_api_key.test", "kind", "team"),
					resource.TestCheckResourceAttrPair("rootly_api_key.test", "group_id", "rootly_team.test", "id"),
				),
			},
		},
	})
}

func testAccResourceApiKeyTeamScopedConfig(name, expiresAt string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}

resource "rootly_api_key" "test" {
	name       = "%s"
	expires_at = "%s"
	kind       = "team"
	group_id   = rootly_team.test.id
}
`, name, name, expiresAt)
}

// group_id is only accepted on create: the API never returns it and rejects it
// on update, so the resource must send it on create only and force replacement
// when it changes.
func TestResourceApiKeyGroupIdIsCreateOnly(t *testing.T) {
	groupId := resourceApiKey().Schema["group_id"]

	if groupId == nil {
		t.Fatal("expected rootly_api_key to expose a group_id attribute")
	}
	if !groupId.Optional {
		t.Error("expected group_id to be optional")
	}
	if !groupId.ForceNew {
		t.Error("expected group_id to force replacement")
	}
}
