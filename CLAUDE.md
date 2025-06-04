# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

### Build and Development
- `make build` - Complete build process: generates code from OpenAPI schema, compiles provider, and regenerates docs
- `make generate` - Downloads Rootly OpenAPI schema and auto-generates client code and provider resources
- `make docs` - Regenerates Terraform documentation from provider schemas
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests (requires `TF_ACC=1` environment variable)
- `go build -o terraform-provider-rootly` - Quick local build without code generation

### Local Testing Setup
```bash
# Build and install locally for testing
make build
mkdir -p ~/.terraform.d/plugins/terraform.local/local/rootly/1.0.0/darwin_arm64/
cp terraform-provider-rootly ~/.terraform.d/plugins/terraform.local/local/rootly/1.0.0/darwin_arm64/terraform-provider-rootly_v1.0.0
```

Configure `~/.terraform.rc` for local testing:
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

## Architecture Overview

### Code Generation System
This provider is heavily auto-generated from Rootly's OpenAPI schema:

1. **Schema Download**: `make generate` fetches the latest OpenAPI spec from Rootly's S3 bucket
2. **Schema Processing**: `tools/clean-swagger.js` cleans the schema, `tools/generate.js` orchestrates generation
3. **Client Generation**: Uses `oapi-codegen` to generate Go client code in `schema/` directory
4. **Provider Generation**: Auto-generates data sources and resources in `provider/` directory
5. **Documentation**: Auto-generates Terraform docs in `docs/` directory

### Key Generation Templates
- `tools/generate-provider-tpl.js` - Main provider configuration
- `tools/generate-resource-tpl.js` - Terraform resource implementations
- `tools/generate-data-source-tpl.js` - Terraform data source implementations
- `tools/generate-client-tpl.js` - Client method implementations
- `tools/generate-workflow-tpl.js` - Workflow-specific resources
- `tools/generate-tasks.js` - Workflow task resources

### Resource Exclusion
Resources can be excluded from generation by adding them to the `excluded` object in `tools/generate.js`. This is useful for resources that need manual implementation or aren't ready for Terraform.

### File Structure Patterns
- **Client Layer**: `client/*.go` - HTTP client implementations for each API resource
- **Provider Layer**: `provider/resource_*.go` and `provider/data_source_*.go` - Terraform resource/data source implementations
- **Schema Layer**: `schema/*.gen.go` - Auto-generated Go structs from OpenAPI schema
- **Tests**: `provider/*_test.go` - Acceptance tests for resources

### Authentication & Configuration
Provider supports:
- `api_host` - Defaults to `https://api.rootly.com`, configurable via `ROOTLY_API_URL`
- `api_token` - Required, configurable via `ROOTLY_API_TOKEN`

### Workflow Tasks Architecture
Workflow tasks are special resources with dynamic generation based on OpenAPI schema task definitions. They follow a pattern of `workflow_task_<action>_<integration>` and are automatically generated from the API schema.

## Development Workflow

1. **Adding New Resources**: Most resources are auto-generated. Add exclusions in `tools/generate.js` only if manual implementation is needed.
2. **Schema Updates**: Run `make generate` to pull latest OpenAPI schema and regenerate all code.
3. **Documentation Updates**: Run `make docs` after any provider schema changes.
4. **Testing**: Always run `make testacc` before submitting changes. Tests require valid Rootly API credentials.
5. **Local Testing**: Use the local installation process above to test provider changes against real Terraform configurations.

## Important Notes

- **Generated Files**: Files marked with "DO NOT MODIFY" headers are auto-generated. Changes should be made to templates in `tools/` directory.
- **Schema Configuration**: `schema/oapi-config.yml` controls OpenAPI code generation, including schema exclusions.
- **Release Process**: New releases are triggered by Git tags and automatically published to Terraform Registry.
- **Node.js Dependency**: Code generation requires Node.js for JavaScript-based template processing.