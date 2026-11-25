package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rootlyhq/terraform-provider-rootly/v5/tools"
)

func TestCatalogPropertyDescriptionCanBeCleared(t *testing.T) {
	mock := newCatalogPropertyDescriptionAPI(t)
	defer mock.server.Close()
	address := "rootly_catalog_property.test"
	config := func(name, description string) string {
		return catalogPropertyDescriptionConfig(mock.server.URL, name, description)
	}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config("Tier", `"Use T0 to T3"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(address, "description", "Use T0 to T3"),
					mock.check(0, nil),
				),
			},
			{Config: config("Tier", `"Use T0 to T3"`), PlanOnly: true},
			{
				Config: config("Tier renamed", "omitted"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(address, "description", "Use T0 to T3"),
					mock.check(1, nil),
				),
			},
			{
				Config: config("Tier renamed", `""`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(address, "description", ""),
					mock.check(2, tools.String("")),
				),
			},
			{Config: config("Tier renamed", `""`), PlanOnly: true},
			{Config: config("Tier renamed", "omitted"), PlanOnly: true},
		},
	})
}

func catalogPropertyDescriptionConfig(host, name, description string) string {
	var config strings.Builder
	fmt.Fprintf(&config, "provider \"rootly\" {\n api_host = %q\n api_token = \"local-test-token\"\n}\n", host)
	fmt.Fprintf(&config, "resource \"rootly_catalog_property\" \"test\" {\n catalog_type = \"service\"\n name = %q\n", name)
	if description != "omitted" {
		fmt.Fprintf(&config, " description = %s\n", description)
	}
	config.WriteString("}\n")
	return config.String()
}

type catalogPropertyDescriptionAPI struct {
	server     *httptest.Server
	mu         sync.Mutex
	attributes map[string]interface{}
	lastUpdate map[string]interface{}
	updates    int
}

func newCatalogPropertyDescriptionAPI(t *testing.T) *catalogPropertyDescriptionAPI {
	t.Helper()
	mock := &catalogPropertyDescriptionAPI{attributes: map[string]interface{}{
		"kind":         "text",
		"catalog_type": "service",
		"multiple":     false,
		"required":     false,
	}}
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/services/properties":
		case r.URL.Path == "/v1/catalog_properties/property-id":
		default:
			t.Errorf("unexpected API request: %s %s", r.Method, r.URL.Path)
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
				// The API stores a blank description as null.
				if field == "description" && value == "" {
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
			"id": "property-id", "type": "catalog_properties", "attributes": mock.attributes,
		}})
	}))
	return mock
}

func (mock *catalogPropertyDescriptionAPI) check(updates int, sent *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		mock.mu.Lock()
		defer mock.mu.Unlock()
		if mock.updates != updates {
			return fmt.Errorf("got %d updates, expected %d", mock.updates, updates)
		}
		if updates == 0 {
			return nil
		}
		got, present := mock.lastUpdate["description"]
		if sent == nil {
			if present {
				return fmt.Errorf("update unexpectedly sent description = %#v", got)
			}
			return nil
		}
		if !present || got != *sent {
			return fmt.Errorf("update description = %#v (present: %t), expected %q", got, present, *sent)
		}
		return nil
	}
}
