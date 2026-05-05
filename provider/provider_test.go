package provider

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/meta"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"rootly": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

var testAccSharedClient *rootly.ClientWithResponses

func init() {
	var apiHost string
	if v := os.Getenv("ROOTLY_API_URL"); v != "" {
		apiHost = v
	} else {
		apiHost = "https://api.rootly.com"
	}

	var apiToken string
	if v := os.Getenv("ROOTLY_API_TOKEN"); v != "" {
		apiToken = v
	}

	legacyClient, err := client.NewClient(apiHost, apiToken, RootlyUserAgent(meta.GetVersion()))
	if err != nil {
		log.Fatalln(err.Error())
	}

	sharedClient, err := rootly.NewClientWithResponses(
		apiHost,
		// Piggyback on the legacy client's HTTP client. Inherits the same headers, authentication, and retry logic.
		rootly.WithHTTPClient(legacyClient),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	testAccSharedClient = sharedClient
}
