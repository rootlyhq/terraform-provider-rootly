# Rootly Terraform Provider Migration Script

This migration script helps users migrate from deprecated Rootly resources to their newer equivalents. It's designed to be generic and scalable for various migration scenarios.

## Overview

The script:
1. Fetches existing resources from your Rootly organization via the API
2. Converts deprecated resources to newer Terraform resource formats
3. Generates Terraform resource configurations
4. Provides import statements/blocks to import existing resources into Terraform state

## Supported Migrations

Currently supported migration types:

### `alert_routing_rules_to_alert_routes`
Migrates from deprecated `rootly_alert_routing_rule` resources to the new `rootly_alert_route` resources.

## Usage

### Basic Usage

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main alert_routing_rules_to_alert_routes > ./alert_routes.tf
```

### With Import Blocks

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -import=block alert_routing_rules_to_alert_routes > ./alert_routes.tf
```

### With Custom API Configuration

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -api-host=https://api.rootly.com -api-token=your-token alert_routing_rules_to_alert_routes > ./alert_routes.tf
```

## Command Line Options

- `migration_type` (required): Specifies which migration to perform. Currently supports:
  - `alert_routing_rules_to_alert_routes`: Migrate from alert routing rules to alert routes
- `-import`: Import statement format. Options: `statement` (default) or `block`
- `-api-host`: Rootly API host (defaults to `https://api.rootly.com` or `ROOTLY_API_URL` env var)
- `-api-token`: Rootly API token (defaults to `ROOTLY_API_TOKEN` env var)

## Environment Variables

- `ROOTLY_API_TOKEN`: Your Rootly API token
- `ROOTLY_API_URL`: Custom Rootly API URL (optional)

## Migration Process

### Step 1: Generate Migration Files

Run the migration script to generate your new `alert_routes.tf` file:

```bash
export ROOTLY_API_TOKEN="your-api-token"
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main alert_routing_rules_to_alert_routes > ./alert_routes.tf
```

### Step 2: Review Generated Configuration

Review the generated `alert_routes.tf` file to ensure all your alert routes have been converted correctly.

### Step 3: Import Resources

Run terraform plan to verify the import operations:

```bash
terraform plan
```

Apply the import operations:

```bash
terraform apply
```

### Step 4: Clean Up Import Statements

Remove the import blocks/statements from your `alert_routes.tf` file as they are no longer needed.

### Step 5: Remove Deprecated Resources

1. Remove all deprecated `rootly_alert_routing_rule` resources from your existing Terraform configuration files
2. Remove the deprecated resources from Terraform state:

```bash
terraform state rm rootly_alert_routing_rule.my_old_rule
terraform state rm rootly_alert_routing_rule.another_old_rule
# ... repeat for each deprecated resource
```

The alert routes are now managed by Terraform using the new `rootly_alert_route` resource type.

## Example Output

The script generates output similar to this:

```hcl
resource "rootly_alert_route" "my_alert_route" {
  name               = "My Alert Route"
  enabled            = true
  alerts_source_ids  = ["source-id-123"]
  owning_team_ids    = ["team-id-456"]

  rules {
    name = "My Rule"
    
    destinations {
      target_type = "Service"
      target_id   = "service-id-789"
    }
    
    condition_groups {
      conditions {
        property_field_type            = "attribute"
        property_field_name            = "summary"
        property_field_condition_type  = "contains"
        property_field_value           = "error"
      }
    }
  }
}

# Import statements
# terraform import rootly_alert_route.my_alert_route 'route-id-123'
```

## Conversion Details

For the `alert_routing_rules_to_alert_routes` migration, the script performs:

1. **Resource Migration**: Converts deprecated `rootly_alert_routing_rule` to new `rootly_alert_route` format
2. **Rule Preservation**: Preserves all existing rules with their conditions and destinations
3. **Team Ownership**: Includes owning team IDs if configured
4. **Naming**: Generates safe Terraform resource names from route names
5. **Import Support**: Provides import statements to bring existing resources under Terraform management

## Adding New Migrations

The script is designed to be extensible. To add a new migration type:

1. Add a new `MigrationType` constant in `program.go`
2. Add the new migration type to `AllMigrationTypes`
3. Add a new case in the `generateOutput()` switch statement
4. Implement the handler function in `converter.go`
5. Update documentation

## Requirements

- Go 1.21 or later
- Valid Rootly API token
- Network access to Rootly API

## Troubleshooting

### Authentication Error
Ensure your API token is valid and has sufficient permissions to read alert routing rules.

### Network Issues
Verify you can reach the Rootly API host. If using a custom host, ensure it's correctly configured.

### Resource Name Conflicts
The script sanitizes resource names automatically. Review generated names and adjust if needed.

## Support

For issues with the migration script, please file a GitHub issue in the terraform-provider-rootly repository.