# Rootly Terraform Provider Migration Script

This migration script helps users import existing Rootly alert routes into their Terraform configuration as `rootly_alert_route` resources.

## Overview

The script:
1. Fetches all alert routes from your Rootly organization via the API
2. Converts each alert route to Terraform `rootly_alert_route` resource format
3. Generates Terraform resource configurations
4. Provides import statements/blocks to import existing resources into Terraform state

## Usage

### Basic Usage

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main alert_routes > ./alert_routes.tf
```

### With Import Blocks

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -import=block alert_routes > ./alert_routes.tf
```

### With Custom API Configuration

```bash
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -api-host=https://api.rootly.com -api-token=your-token alert_routes > ./alert_routes.tf
```

## Command Line Options

- `object_type` (required): Currently supports `alert_routes`
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
go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main alert_routes > ./alert_routes.tf
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

Remove the import blocks/statements from your `alert_routes.tf` file as they are no longer needed. The alert routes are now managed by Terraform.

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

The script performs the following conversions:

1. **Direct Mapping**: Maps alert route API responses to Terraform resource format
2. **Rule Preservation**: Preserves all existing rules with their conditions and destinations
3. **Team Ownership**: Includes owning team IDs if configured
4. **Naming**: Generates safe Terraform resource names from alert route names

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