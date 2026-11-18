package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestStatusPageSamlClearableFields(t *testing.T) {
	certificate, _ := testSamlCertificate(t)
	cases := []struct {
		kind       string
		collection string
		attributes map[string]string
		fields     []string
	}{
		{
			kind:       "status_page",
			collection: "status-pages",
			attributes: map[string]string{"title": "Acme"},
			fields:     []string{"saml_idp_entity_id", "saml_configuration_id"},
		},
		{
			kind:       "status_page_saml_configuration",
			collection: "status-page-saml-configurations",
			attributes: map[string]string{
				"name":                "Okta IdP",
				"idp_sso_service_url": "https://idp.example.com/sso",
				"idp_entity_id":       "https://idp.example.com/entity",
				"idp_cert":            certificate,
			},
			fields: []string{"idp_slo_service_url"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			mock := newSamlClearableAPI(t, tc.collection)
			defer mock.server.Close()
			address := "rootly_" + tc.kind + ".test"
			config := func(rename, value string) string {
				return samlClearableConfig(mock.server.URL, tc.kind, tc.attributes, rename, tc.fields, value)
			}
			check := func(updates int, sent, saved string) resource.TestCheckFunc {
				checks := []resource.TestCheckFunc{mock.check(updates, tc.fields, sent, saved)}
				for _, field := range tc.fields {
					checks = append(checks, resource.TestCheckResourceAttr(address, field, samlClearableValue(saved, field)))
				}
				return resource.ComposeTestCheckFunc(checks...)
			}
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: providerFactories,
				Steps: []resource.TestStep{
					{Config: config("", "set-"), Check: check(0, "set-", "set-")},
					{Config: config("", "set-"), PlanOnly: true},
					{Config: config("renamed", "omitted"), Check: check(1, "omitted", "set-")},
					{Config: config("renamed", "omitted"), PlanOnly: true},
					{Config: config("renamed", ""), Check: check(2, "", "")},
					{Config: config("renamed", ""), PlanOnly: true},
					{Config: config("renamed", "omitted"), PlanOnly: true},
					{Config: config("renamed", "restored-"), Check: check(3, "restored-", "restored-")},
					{Config: config("renamed", "restored-"), PlanOnly: true},
				},
			})
		})
	}
}

func samlClearableValue(prefix, field string) string {
	if prefix == "" {
		return ""
	}
	return prefix + field
}

func samlClearableConfig(host, kind string, attributes map[string]string, rename string, fields []string, value string) string {
	var config strings.Builder
	fmt.Fprintf(&config, "provider \"rootly\" {\n api_host = %q\n api_token = \"local-test-token\"\n}\n", host)
	fmt.Fprintf(&config, "resource \"rootly_%s\" \"test\" {\n", kind)
	names := make([]string, 0, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		attribute := attributes[name]
		if rename != "" && (name == "title" || name == "name") {
			attribute += " " + rename
		}
		fmt.Fprintf(&config, " %s = %q\n", name, attribute)
	}
	for _, field := range fields {
		switch value {
		case "omitted":
		case "":
			fmt.Fprintf(&config, " %s = \"\"\n", field)
		default:
			fmt.Fprintf(&config, " %s = %q\n", field, value+field)
		}
	}
	config.WriteString("}\n")
	return config.String()
}

type samlClearableAPI struct {
	server     *httptest.Server
	mu         sync.Mutex
	attributes map[string]interface{}
	lastUpdate map[string]interface{}
	updates    int
}

func newSamlClearableAPI(t *testing.T, collection string) *samlClearableAPI {
	t.Helper()
	mock := &samlClearableAPI{attributes: map[string]interface{}{}}
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		if r.URL.Path != "/v1/"+collection && r.URL.Path != "/v1/"+collection+"/test-id" {
			t.Errorf("unexpected API path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var payload struct {
				Data struct {
					Attributes map[string]interface{} `json:"attributes"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Method == http.MethodPut {
				mock.updates++
				mock.lastUpdate = payload.Data.Attributes
			}
			for field, value := range payload.Data.Attributes {
				if value == "" {
					value = nil
				}
				mock.attributes[field] = value
			}
		case http.MethodGet:
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			t.Errorf("unexpected API method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": "test-id", "type": strings.ReplaceAll(collection, "-", "_"), "attributes": mock.attributes,
		}})
	}))
	return mock
}

func (mock *samlClearableAPI) check(updates int, fields []string, sent, saved string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		if mock.updates != updates {
			return fmt.Errorf("got %d updates, expected %d", mock.updates, updates)
		}
		for _, field := range fields {
			gotSaved, _ := mock.attributes[field].(string)
			if want := samlClearableValue(saved, field); gotSaved != want {
				return fmt.Errorf("remote %s = %q, expected %q", field, gotSaved, want)
			}
			if updates == 0 {
				continue
			}
			got, present := mock.lastUpdate[field]
			if sent == "omitted" {
				if present {
					return fmt.Errorf("update unexpectedly sent %s = %#v", field, got)
				}
				continue
			}
			if want := samlClearableValue(sent, field); !present || got != want {
				return fmt.Errorf("update %s = %#v (present: %t), expected %q", field, got, present, want)
			}
		}
		return nil
	}
}
