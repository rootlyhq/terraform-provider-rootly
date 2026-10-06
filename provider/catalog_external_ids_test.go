package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var catalogExternalIDFields = []string{
	"backstage_id", "external_id", "pagerduty_id", "opsgenie_id",
	"opsgenie_team_id", "cortex_id", "opslevel_id",
}

func TestCatalogExternalIDTerraformRoundTrip(t *testing.T) {
	for _, kind := range []string{"service", "functionality"} {
		t.Run(kind, func(t *testing.T) {
			name := acctest.RandomWithPrefix("tf-" + kind)
			mock := newCatalogExternalIDAPI(t, kind)
			defer mock.server.Close()
			address := "rootly_" + kind + ".test"
			fields := catalogManagedExternalIDFields(kind)
			config := func(name, value string) string {
				return catalogExternalIDConfig(mock.server.URL, kind, name, fields, value)
			}
			check := func(updates int, sent, saved string) resource.TestCheckFunc {
				checks := []resource.TestCheckFunc{mock.check(updates, fields, sent, saved)}
				for _, field := range fields {
					value := saved
					if saved != "" {
						value += field
					}
					checks = append(checks, resource.TestCheckResourceAttr(address, field, value))
				}
				return resource.ComposeTestCheckFunc(checks...)
			}
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: providerFactories,
				Steps: []resource.TestStep{
					{Config: config(name, "omitted"), Check: check(0, "omitted", "remote-")},
					{Config: config(name, "omitted"), PlanOnly: true},
					{Config: config(name+"-renamed", "omitted"), Check: check(1, "omitted", "remote-")},
					{ResourceName: address, ImportState: true, ImportStateVerify: true},
					{Config: config(name+"-renamed", "omitted"), PlanOnly: true},
					{Config: config(name+"-renamed", "replacement-"), Check: check(2, "replacement-", "replacement-")},
					{Config: config(name+"-renamed", "replacement-"), PlanOnly: true},
					{Config: config(name+"-unmanaged", "omitted"), Check: check(3, "omitted", "replacement-")},
					{Config: config(name+"-unmanaged", "omitted"), PlanOnly: true},
					{Config: config(name+"-null", "null"), Check: check(4, "omitted", "replacement-")},
					{Config: config(name+"-null", "null"), PlanOnly: true},
					{Config: config(name+"-null", ""), Check: check(5, "", "")},
					{Config: config(name+"-null", ""), PlanOnly: true},
					{Config: config(name+"-null", "omitted"), PlanOnly: true},
					{Config: config(name+"-null", "restored-"), Check: check(6, "restored-", "restored-")},
					{Config: config(name+"-null", "restored-"), PlanOnly: true},
				},
			})
		})
	}
}

func TestCatalogExternalIDExistingServiceState(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-service-upgrade")
	mock := newCatalogExternalIDAPI(t, "service")
	defer mock.server.Close()
	fields := catalogManagedExternalIDFields("service")
	oldFactories := map[string]func() (*schema.Provider, error){
		"rootly": func() (*schema.Provider, error) {
			p := New("dev")()
			old := resourceService()
			for _, field := range fields {
				old.Schema[field].Computed = false
			}
			p.ResourcesMap["rootly_service"] = old
			return p, nil
		},
	}
	resource.UnitTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProviderFactories: oldFactories,
				Config:            catalogExternalIDConfig(mock.server.URL, "service", name, fields, "remote-"),
			},
			{
				ProviderFactories: providerFactories,
				Config:            catalogExternalIDConfig(mock.server.URL, "service", name, fields, "omitted"),
				PlanOnly:          true,
			},
			{
				ProviderFactories: providerFactories,
				Config:            catalogExternalIDConfig(mock.server.URL, "service", name+"-upgraded", fields, "omitted"),
				Check:             mock.check(1, fields, "omitted", "remote-"),
			},
			{
				ProviderFactories: providerFactories,
				Config:            catalogExternalIDConfig(mock.server.URL, "service", name+"-upgraded", fields, "omitted"),
				PlanOnly:          true,
			},
		},
	})
}

func catalogManagedExternalIDFields(kind string) []string {
	if kind == "service" {
		return []string{"backstage_id", "external_id", "pagerduty_id", "opsgenie_id", "cortex_id"}
	}
	return []string{"backstage_id", "external_id", "pagerduty_id", "opsgenie_id", "opsgenie_team_id", "cortex_id"}
}

func catalogExternalIDConfig(host, kind, name string, fields []string, value string) string {
	var config strings.Builder
	fmt.Fprintf(&config, "provider \"rootly\" {\n api_host = %q\n api_token = \"local-test-token\"\n}\n", host)
	fmt.Fprintf(&config, "resource \"rootly_%s\" \"test\" {\n name = %q\n", kind, name)
	for _, field := range fields {
		switch value {
		case "omitted":
		case "null":
			fmt.Fprintf(&config, " %s = null\n", field)
		case "":
			fmt.Fprintf(&config, " %s = \"\"\n", field)
		default:
			fmt.Fprintf(&config, " %s = %q\n", field, value+field)
		}
	}
	config.WriteString("}\n")
	return config.String()
}

type catalogExternalIDAPI struct {
	server     *httptest.Server
	mu         sync.Mutex
	attributes map[string]interface{}
	lastUpdate map[string]interface{}
	updates    int
}

func newCatalogExternalIDAPI(t *testing.T, kind string) *catalogExternalIDAPI {
	t.Helper()
	collection := kind + "s"
	if kind == "functionality" {
		collection = "functionalities"
	}
	mock := &catalogExternalIDAPI{attributes: map[string]interface{}{}}
	for _, field := range catalogExternalIDFields {
		mock.attributes[field] = "remote-" + field
	}
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		if r.URL.Path != "/v1/"+collection && r.URL.Path != "/v1/"+collection+"/catalog-id" {
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
				// The API normalizes intentionally cleared IDs to null on read.
				if strings.HasSuffix(field, "_id") && value == "" {
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
			"id": "catalog-id", "type": collection, "attributes": mock.attributes,
		}})
	}))
	return mock
}

func (mock *catalogExternalIDAPI) check(updates int, fields []string, sent, saved string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		if mock.updates != updates {
			return fmt.Errorf("got %d updates, expected %d", mock.updates, updates)
		}
		for _, field := range catalogExternalIDFields {
			managed := false
			for _, candidate := range fields {
				managed = managed || field == candidate
			}
			wantSaved := "remote-" + field
			if managed {
				wantSaved = saved
				if saved != "" {
					wantSaved += field
				}
			}
			gotSaved, _ := mock.attributes[field].(string)
			if gotSaved != wantSaved {
				return fmt.Errorf("remote %s = %q, expected %q", field, gotSaved, wantSaved)
			}
			if updates == 0 {
				continue
			}
			got, present := mock.lastUpdate[field]
			if sent == "omitted" || !managed {
				if present {
					return fmt.Errorf("update unexpectedly sent %s = %#v", field, got)
				}
			} else {
				want := sent
				if sent != "" {
					want += field
				}
				if !present || got != want {
					return fmt.Errorf("update %s = %#v (present: %t), expected %q", field, got, present, want)
				}
			}
		}
		return nil
	}
}
