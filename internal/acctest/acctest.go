package acctest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	sdkv2_provider "github.com/rootlyhq/terraform-provider-rootly/v2/provider"
)

func PreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"rootly": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkServer, err := tf5to6server.UpgradeServer(
			ctx,
			sdkv2_provider.New("dev")().GRPCProvider,
		)
		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			func() tfprotov6.ProviderServer {
				return upgradedSdkServer
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
