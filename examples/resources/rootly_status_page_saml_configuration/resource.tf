resource "rootly_status_page_saml_configuration" "okta" {
  name                = "Okta"
  idp_sso_service_url = "https://acme.okta.com/app/acme_status/sso/saml"
  idp_entity_id       = "http://www.okta.com/exk1example"
  idp_cert            = file("${path.module}/okta.pem")
}

resource "rootly_status_page" "internal" {
  title                   = "Acme Internal Status"
  authentication_method   = "saml"
  saml_configuration_id   = rootly_status_page_saml_configuration.okta.id
  saml_access_rules_match = "any"

  saml_access_rules {
    claim    = "groups"
    operator = "is_any_of"
    values   = ["engineering", "support"]
  }

  saml_access_rules {
    operator = "email_domain_is_any_of"
    values   = ["acme.com"]
  }
}
