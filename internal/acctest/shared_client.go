package acctest

import (
	"log"
	"os"

	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/meta"
	sdkv2_provider "github.com/rootlyhq/terraform-provider-rootly/v2/provider"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

var (
	TestApiHost  string
	TestApiToken string

	SharedClient *rootly.ClientWithResponses
)

func init() {
	if v := os.Getenv("ROOTLY_API_URL"); v != "" {
		TestApiHost = v
	} else {
		TestApiHost = "https://api.rootly.com"
	}

	if v := os.Getenv("ROOTLY_API_TOKEN"); v != "" {
		TestApiToken = v
	}

	legacyClient, err := client.NewClient(TestApiHost, TestApiToken, sdkv2_provider.RootlyUserAgent(meta.GetVersion()))
	if err != nil {
		log.Fatalf("Unable to create Rootly client: %v", err)
	}

	client, err := rootly.NewClientWithResponses(
		TestApiHost,
		// Piggyback on the legacy client's HTTP client. Inherits the same headers, authentication, and retry logic.
		rootly.WithHTTPClient(legacyClient),
	)
	if err != nil {
		log.Fatalf("Unable to create Rootly client: %v", err)
	}

	SharedClient = client
}
