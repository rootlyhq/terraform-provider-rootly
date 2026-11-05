package provider

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestResourceStatusPageSamlConfigurationSchema(t *testing.T) {
	t.Parallel()

	s := resourceStatusPageSamlConfiguration().Schema

	require.True(t, s["idp_cert"].Sensitive)
	require.NotNil(t, s["idp_cert"].DiffSuppressFunc)
	for _, name := range []string{"idp_cert_fingerprint", "sp_entity_id", "acs_url", "metadata_url", "status_pages_count"} {
		require.True(t, s[name].Computed, name)
		require.False(t, s[name].Optional, name)
	}

	page := resourceStatusPage().Schema
	for _, name := range []string{"saml_configuration_id", "saml_idp_entity_id"} {
		require.True(t, page[name].Optional, name)
		require.True(t, page[name].Computed, name)
	}
	require.True(t, page["saml_access_rules"].Optional)
	require.True(t, page["saml_access_rules_match"].Optional)
}

func TestAccResourceStatusPageSamlConfiguration(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-saml-config")
	cert := testAccStatusPageSamlCertificate(t)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStatusPageSamlConfigurationConfig(rName, cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_status_page_saml_configuration.test", "name", rName),
					resource.TestCheckResourceAttrSet("rootly_status_page_saml_configuration.test", "idp_cert_fingerprint"),
					resource.TestCheckResourceAttrSet("rootly_status_page_saml_configuration.test", "acs_url"),
					resource.TestCheckResourceAttrPair("rootly_status_page.test", "saml_configuration_id", "rootly_status_page_saml_configuration.test", "id"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.#", "1"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.0.claim", "groups"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules.0.values.#", "2"),
					resource.TestCheckResourceAttr("rootly_status_page.test", "saml_access_rules_match", "any"),
				),
			},
			{
				Config:   testAccResourceStatusPageSamlConfigurationConfig(rName, cert),
				PlanOnly: true,
			},
		},
	})
}

func testAccStatusPageSamlCertificate(t *testing.T) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: "tf-acc-saml-idp"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(1, 0, 0),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func testAccResourceStatusPageSamlConfigurationConfig(name, cert string) string {
	return fmt.Sprintf(`
resource "rootly_status_page_saml_configuration" "test" {
  name                = "%[1]s"
  idp_sso_service_url = "https://idp.example.com/sso"
  idp_entity_id       = "https://idp.example.com/%[1]s"
  idp_cert            = <<-EOT
%[2]s
EOT
}

resource "rootly_status_page" "test" {
  title                   = "%[1]s"
  authentication_enabled  = true
  authentication_method   = "saml"
  saml_configuration_id   = rootly_status_page_saml_configuration.test.id
  saml_access_rules_match = "any"

  saml_access_rules {
    claim    = "groups"
    operator = "is_any_of"
    values   = ["engineering", "support"]
  }
}
`, name, cert)
}
