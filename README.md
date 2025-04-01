# Rootly Provider

The [Rootly](https://rootly.com/) provider is used to interact with the resources supported by Rootly. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 0.14 or later.

## Usage

Please see the [Terraform Registry documentation](https://registry.terraform.io/providers/rootlyhq/rootly/latest/docs) or [examples/](examples).
}

## Develop

`make build` auto-generates code from Rootly's JSON-API schema, compiles provider code, and regenerates docs.

Exclude API resources from the provider by adding them to the ignore-list in `tools/generate.js`.

## Release

Tag a new release to automatically build and publish to the Terraform Registry.
