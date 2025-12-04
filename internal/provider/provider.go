package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/apiclient"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RootlyProviderModel struct {
	ApiHost  types.String `tfsdk:"api_host"`
	ApiToken types.String `tfsdk:"api_token"`
}

type RootlyProviderData struct {
	Client       *rootly.ClientWithResponses
	LegacyClient *client.Client
}

var _ provider.Provider = &RootlyProvider{}

type RootlyProvider struct {
	version string
}

func (p *RootlyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rootly"
	resp.Version = p.version
}

func (p *RootlyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_host": schema.StringAttribute{
				MarkdownDescription: "The Rootly API host. Defaults to https://api.rootly.com. Can also be sourced from the `ROOTLY_API_URL` environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The Rootly API Token. Generate it from your account at https://rootly.com/account. It must be provided but can also be sourced from the `ROOTLY_API_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *RootlyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RootlyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiHost string
	if !data.ApiHost.IsNull() {
		apiHost = data.ApiHost.ValueString()
	} else if v := os.Getenv("ROOTLY_API_URL"); v != "" {
		apiHost = v
	} else {
		apiHost = "https://api.rootly.com"
	}

	var apiToken string
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	} else if v := os.Getenv("ROOTLY_API_TOKEN"); v != "" {
		apiToken = v
	}

	legacyClient, client, err := apiclient.New(apiHost, apiToken, p.version)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Rootly client", err.Error())
		return
	}

	providerData := &RootlyProviderData{
		Client:       client,
		LegacyClient: legacyClient,
	}

	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

// DataSources implements provider.Provider.
func (p *RootlyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources implements provider.Provider.
func (p *RootlyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RootlyProvider{
			version: version,
		}
	}
}
