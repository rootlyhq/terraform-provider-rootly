# Rootly Provider

The [Rootly](https://rootly.com/) provider is used to interact with the resources supported by Rootly. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 0.14 or later.

## Usage

Please see the [Terraform Registry documentation](https://registry.terraform.io/providers/rootlyhq/rootly/latest/docs) or [examples/](examples).

## Develop

### Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Complete build: generate code, compile provider, regenerate docs |
| `make generate` | Download schema and auto-generate client code and provider resources |
| `make docs` | Regenerate Terraform documentation from provider schemas |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests |
| `make help` | Show all available commands |

### Updating provider

`make build` auto-generates code from Rootly's JSON-API schema, compiles provider code, and regenerates docs.

Some API resources are excluded from code generation if they are in the ignore-lists in `tools/generate.js`. If so, those resources will need to be updated manually.

Tests are often not able to be code generated. If so, add them to the ignore-list in `tools/generate.js` and implement manually.

### Version Management

The project uses [semantic versioning](https://semver.org) with git tags. Use these commands to manage versions:

#### Version Commands
```bash
make version-show      # Show current and next versions
make version-patch     # Bump patch version (1.2.3 → 1.2.4)
make version-minor     # Bump minor version (1.2.3 → 1.3.0)  
make version-major     # Bump major version (1.2.3 → 2.0.0)
```

#### Release Commands
```bash
make release-patch     # Bump patch version + push tag (triggers CI release)
make release-minor     # Bump minor version + push tag (triggers CI release)
make release-major     # Bump major version + push tag (triggers CI release)
```

#### Example Usage
```bash
# Check current version
make version-show
# Current version: v3.1.0
# Next patch: 3.1.1

# Bump patch version
make version-patch
# Creates and pushes v3.1.1 tag

# Or bump and push tag in one step (triggers CI release)
make release-patch
# Bumps to v3.1.1, pushes tag, CI builds and publishes release
```

**Important**: Use the `make version-*` commands instead of manually creating git tags. This ensures version consistency and proper validation.

### Release / Publish to Terraform Registry

Releases are automatically published to Terraform Registry when a new tag is pushed. The version management commands above handle this process:

1. **Version Bumping**: `make version-*` commands create and push git tags
2. **CI Trigger**: Pushed tags trigger CI/GoReleaser workflow 
3. **Automatic Release**: CI builds and publishes releases to GitHub and Terraform Registry
4. **Version Injection**: The correct version is automatically set in the provider binary during CI build

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
