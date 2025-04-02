# Rootly Provider

The [Rootly](https://rootly.com/) provider is used to interact with the resources supported by Rootly. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 0.14 or later.

## Usage

Please see the [Terraform Registry documentation](https://registry.terraform.io/providers/rootlyhq/rootly/latest/docs) or [examples/](examples).
}

## Develop

`make build` auto-generates code from Rootly's JSON-API schema, compiles provider code, and regenerates docs.

Exclude API resources from the provider by adding them to the ignore-list in `tools/generate.js`.

### Local build

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

## Release

Tag a new release to automatically build and publish to the Terraform Registry.
