package provider

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceStatusPageSamlConfigurationIdpCertDiff(t *testing.T) {
	current, currentFingerprint := testSamlCertificate(t)
	rotated, _ := testSamlCertificate(t)

	state := &terraform.InstanceState{
		ID: "c0ffee00-0000-4000-8000-000000000001",
		Attributes: map[string]string{
			"id":                     "c0ffee00-0000-4000-8000-000000000001",
			"name":                   "Okta IdP",
			"idp_sso_service_url":    "https://idp.example.com/sso",
			"idp_slo_service_url":    "",
			"idp_entity_id":          "https://idp.example.com/entity",
			"idp_cert":               "[REDACTED]",
			"idp_cert_fingerprint":   currentFingerprint,
			"name_identifier_format": "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
			"sp_entity_id":           "https://rootly.example.com/sp",
			"acs_url":                "https://rootly.example.com/acs",
			"metadata_url":           "https://rootly.example.com/metadata",
			"status_pages_count":     "1",
		},
	}

	cases := []struct {
		name        string
		certificate string
		wantDiff    bool
	}{
		{name: "unchanged certificate", certificate: current, wantDiff: false},
		{name: "unchanged certificate with surrounding whitespace", certificate: "\n" + current + "\n", wantDiff: false},
		{name: "rotated certificate", certificate: rotated, wantDiff: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":                "Okta IdP",
				"idp_sso_service_url": "https://idp.example.com/sso",
				"idp_entity_id":       "https://idp.example.com/entity",
				"idp_cert":            tc.certificate,
			})

			diff, err := resourceStatusPageSamlConfiguration().Diff(context.Background(), state, config, nil)
			if err != nil {
				t.Fatalf("diff failed: %s", err)
			}

			var attribute *terraform.ResourceAttrDiff
			if diff != nil {
				attribute = diff.Attributes["idp_cert"]
			}
			if tc.wantDiff && (attribute == nil || attribute.New != tc.certificate) {
				t.Fatalf("expected idp_cert to change to the configured certificate, got %#v", attribute)
			}
			if !tc.wantDiff && attribute != nil {
				t.Fatalf("expected no idp_cert diff, got %#v", attribute)
			}
		})
	}
}

func testSamlCertificate(t *testing.T) (string, string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "status-page-saml-test-idp"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	sum := sha1.Sum(der)
	hexPairs := make([]string, len(sum))
	for i, b := range sum {
		hexPairs[i] = fmt.Sprintf("%02X", b)
	}

	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})), strings.Join(hexPairs, ":")
}
