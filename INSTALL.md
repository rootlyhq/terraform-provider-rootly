# Terraform Provider Rootly

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Using the provider

To use the provider, refer to the Terraform [provider API documentation](./docs/index.md) and the [Rootly](https://docs.rootly.com/api-reference) API documentation.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

The provider depends on [rootly-go](https://github.com/rootlyhq/rootly-go) for Go types of the Rootly API. These types are auto-generated from Swagger/OpenAPI. See the rootly-go [README](https://github.com/rootlyhq/rootly-go#readme) for more information.