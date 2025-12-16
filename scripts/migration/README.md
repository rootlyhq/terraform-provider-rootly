# Rootly Terraform Provider Migration Script

This migration script helps users migrate from deprecated Rootly resources to their newer equivalents. It's designed to be generic and extensible for various migration scenarios.

## Overview

The migration script automates the process of migrating deprecated Terraform resources to new ones by:

1. **Fetching**: Retrieves existing resources from your Rootly organization via the API using pagination
2. **Converting**: Transforms deprecated resource structures to new Terraform resource formats
3. **Generating**: Creates complete Terraform resource configurations with proper syntax
4. **Importing**: Provides import statements/blocks to bring existing resources under Terraform management
5. **Instructing**: Includes step-by-step instructions for completing the migration process

## Supported Migrations

Currently supported migration types:

### `alert_routing_rules_to_alert_routes`
Migrates from deprecated `rootly_alert_routing_rule` resources to the new `rootly_alert_route` resources. This migration:
- Fetches all alert routes from your Rootly organization 
- Preserves all rule configurations, conditions, and destinations
- Maintains team ownership and alert source associations
- Generates proper Terraform import statements for state management

## Usage

The script follows this command pattern:
```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master [FLAGS] <migration_type> > output_file.tf
```

### Basic Usage with Import Commands

```bash
# Set your API token first
export ROOTLY_API_TOKEN="your-api-token"

# Generate migration file with shell import commands
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master alert_routing_rules_to_alert_routes > alert_routes.tf
```

### Using Import Blocks (Terraform 1.5+)

```bash
# Generate migration file with Terraform import blocks
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master -import=block alert_routing_rules_to_alert_routes > alert_routes.tf
```

### With Custom API Configuration

```bash
# Override default API settings
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master -api-host=https://api.rootly.com -api-token=your-token alert_routing_rules_to_alert_routes > alert_routes.tf
```

## Command Line Options

**Migration Type** (required, last argument):
- `alert_routing_rules_to_alert_routes`: Migrate from deprecated alert routing rules to alert routes

**Flags**:
- `-import`: Import statement format - `statement` (default) or `block`
  - `statement`: Generates shell commands like `terraform import resource.name id`
  - `block`: Generates import blocks for use with `terraform plan/apply` workflow
- `-api-host`: Rootly API host URL (default: `https://api.rootly.com` or `ROOTLY_API_URL` env var)
- `-api-token`: Rootly API authentication token (required, or use `ROOTLY_API_TOKEN` env var)

## Environment Variables

- `ROOTLY_API_TOKEN`: Your Rootly API token
- `ROOTLY_API_URL`: Custom Rootly API URL (optional)

## Migration Process

### Step 1: Generate Migration Files

Set your API token and run the migration script:

```bash
export ROOTLY_API_TOKEN="your-api-token"
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master alert_routing_rules_to_alert_routes > alert_routes.tf
```

### Step 2: Review Generated Configuration

Open and review the generated `alert_routes.tf` file to ensure:
- All your alert routes were fetched correctly
- Resource names are appropriate 
- Rule configurations, conditions, and destinations are preserved
- Import statements/blocks are present at the bottom

### Step 3: Import Resources into Terraform State

The process depends on which import format you used:

#### If using import statements (default):
```bash
# Run the import commands from the generated file
terraform import rootly_alert_route.my_route 'route-id-123'
terraform import rootly_alert_route.another_route 'route-id-456'
# ... (run all import commands from the file)

# Verify resources are imported correctly
terraform plan
# Should show "No changes" if imports were successful
```

#### If using import blocks (`-import=block`):
```bash
# Plan shows import operations
terraform plan

# Apply the import operations  
terraform apply

# Verify import success
terraform plan
# Should show "No changes" after imports complete
```

### Step 4: Clean Up Import Statements

Remove the import statements or import blocks from your `alert_routes.tf` file as they are no longer needed after the resources are imported.

### Step 5: Remove Deprecated Resources

1. **Remove from Terraform state**: Remove all deprecated `rootly_alert_routing_rule` resources from Terraform state:

```bash
terraform state rm rootly_alert_routing_rule.my_old_rule
terraform state rm rootly_alert_routing_rule.another_old_rule
# Repeat for each deprecated resource
```

2. **Remove from configuration**: Remove all deprecated `rootly_alert_routing_rule` resource blocks from your existing Terraform configuration files.

3. **Final verification**: Run `terraform plan` to ensure no unexpected changes.

Your alert routes are now managed by Terraform using the new `rootly_alert_route` resource type.

## Example Output

The script generates output similar to this:

```hcl
resource "rootly_alert_route" "my_alert_route" {
  name               = "My Alert Route"
  enabled            = true
  alerts_source_ids  = ["source-id-123"]
  owning_team_ids    = ["team-id-456"]

  rules {
    name          = "My Rule"
    position      = 1
    fallback_rule = false
    
    destinations {
      target_type = "Service"
      target_id   = "service-id-789"
    }
    
    condition_groups {
      position = 0
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

# Instructions:
# 1. Run the terraform import commands above to import existing resources into state
# 2. Run 'terraform plan' to verify the resources are properly imported
# 3. Remove the import statements above from this file once imports are complete
# 4. Remove deprecated 'rootly_alert_routing_rule' resources from your configuration
# 5. Clean up old resources from Terraform state using 'terraform state rm'
```

## Conversion Details

For the `alert_routing_rules_to_alert_routes` migration, the script performs:

1. **API Fetching**: Retrieves all alert routes from your Rootly organization using paginated requests
2. **Resource Migration**: Converts API response to `rootly_alert_route` Terraform resource format  
3. **Field Mapping**: Maps all fields including:
   - Basic route properties (name, enabled, alerts_source_ids, owning_team_ids)
   - Rule details (name, position, fallback_rule status)
   - Condition groups with positioning
   - Individual conditions with all property types and values
   - Destinations with target types and IDs
4. **Name Sanitization**: Generates valid Terraform resource names from route names
5. **Import Generation**: Creates appropriate import statements/blocks for Terraform state management
6. **Instructions**: Includes detailed step-by-step migration instructions

## Adding New Migrations

The script is designed to be extensible. To add a new migration type:

1. Add a new `MigrationType` constant in `cli.go`
2. Add the new migration type to `AllMigrationTypes` slice
3. Add a new case in the `generateOutput()` switch statement  
4. Implement the handler function in `migrators/` directory
5. Update help text and documentation
6. Add appropriate API models and Terraform templates

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