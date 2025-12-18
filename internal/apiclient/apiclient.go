package apiclient

import (
	"fmt"

	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	sdkv2_provider "github.com/rootlyhq/terraform-provider-rootly/v2/provider"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func New(apiHost string, apiToken string, version string) (*client.Client, *rootly.ClientWithResponses, error) {
	legacyClient, err := client.NewClient(apiHost, apiToken, sdkv2_provider.RootlyUserAgent(version))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Rootly client: %v", err)
	}

	client, err := rootly.NewClientWithResponses(
		apiHost,
		// Piggyback on the legacy client's HTTP client. Inherits the same headers, authentication, and retry logic.
		rootly.WithHTTPClient(legacyClient),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create Rootly client: %v", err)
	}

	return legacyClient, client, nil
}
