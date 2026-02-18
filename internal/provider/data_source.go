package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type baseDataSource struct {
	client       *rootly.ClientWithResponses
	legacyClient *client.Client
}

func (d *baseDataSource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*RootlyProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *RootlyProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = providerData.Client
	d.legacyClient = providerData.LegacyClient
}
