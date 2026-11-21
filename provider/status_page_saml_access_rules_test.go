package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStatusPageSamlAccessRulesRemovingFirstRuleKeepsNoStaleClaim(t *testing.T) {
	mock := newSamlClearableAPI(t, "status-pages")
	defer mock.server.Close()

	groupsRule := `
  saml_access_rules {
    claim    = "groups"
    operator = "is_any_of"
    values   = ["001A"]
  }`
	emailRule := `
  saml_access_rules {
    operator = "email_domain_is_any_of"
    values   = ["acme.com"]
  }`
	config := func(rules ...string) string {
		return fmt.Sprintf(`provider "rootly" {
  api_host  = %q
  api_token = "local-test-token"
}
resource "rootly_status_page" "test" {
  title = "Acme"
%s
}
`, mock.server.URL, strings.Join(rules, "\n"))
	}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config(groupsRule, emailRule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.#", "2"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.1.claim", ""),
				),
			},
			{Config: config(groupsRule, emailRule), PlanOnly: true},
			{
				Config: config(emailRule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.#", "1"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.0.operator", "email_domain_is_any_of"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.0.claim", ""),
					mock.checkSentAccessRuleClaims(""),
				),
			},
			{Config: config(emailRule), PlanOnly: true},
		},
	})
}

func (mock *samlClearableAPI) checkSentAccessRuleClaims(want ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		rules, ok := mock.lastUpdate["saml_access_rules"].([]interface{})
		if !ok || len(rules) != len(want) {
			return fmt.Errorf("update sent saml_access_rules = %#v, expected %d rules", mock.lastUpdate["saml_access_rules"], len(want))
		}
		for i, rule := range rules {
			claim, _ := rule.(map[string]interface{})["claim"].(string)
			if claim != want[i] {
				return fmt.Errorf("update sent rule %d claim %q, expected %q", i, claim, want[i])
			}
		}
		return nil
	}
}
