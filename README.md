# Rootly Provider

The [Rootly](https://rootly.com/) provider is used to interact with the resources supported by Rootly. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 0.14 or later.

## Usage

Please see the [Terraform Registry documentation](https://registry.terraform.io/providers/rootlyhq/rootly/latest/docs) or [examples/](examples).

## Develop

### Updating provider

`make build` auto-generates code from Rootly's JSON-API schema, compiles provider code, and regenerates docs.

Some API resources are excluded from code generation if they are in the ignore-lists in `tools/generate.js`. If so, those resources will need to be updated manually.

Tests are often not able to be code generated. If so, add them to the ignore-list in `tools/generate.js` and implement manually.

### Release / Publish to Terraform Registry

Releases are automatically published to Terraform Registry when a new tag is pushed to main. Tags are [semantic versions](https://semver.org). Before tagging a new release, please update the changelog.

### Local build

```
# build local terraform provider, note every time you make a change
go build -o terraform-provider-rootly

```

Configure a local Terraform registry `terraform.local` in `~/.terraform.rc`

```
provider_installation {
  filesystem_mirror {
    path    = "~/.terraform.d/plugins"
  }
  direct {
    exclude = ["terraform.local/*/*"]
  }
}
```

After building, copy the plugin to the local registry. Change the architecture label as necessary.

```
# make the directory
mkdir -p ~/.terraform.d/plugins/terraform.local/local/rootly/1.0.0/darwin_arm64/

# copy the provider
cp terraform-provider-rootly ~/.terraform.d/plugins/terraform.local/local/rootly/1.0.0/darwin_arm64/terraform-provider-rootly_v1.0.0
```

In your Terraform configuration, specify the local plugin and version 1.0.0:

```
terraform {
  required_providers {
    rootly = {
      source = "terraform.local/local/rootly"
      version = "1.0.0"
    }
  }
}
```

Now running `terraform init` will install the local build.
