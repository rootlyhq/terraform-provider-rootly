package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration/migrators"
)

const (
	MigrationTypeAlertRoutingRulesToAlertRoutes migrators.MigrationType = "alert_routing_rules_to_alert_routes"
)

var AllMigrationTypes = []migrators.MigrationType{
	MigrationTypeAlertRoutingRulesToAlertRoutes,
}

func ToMigrationType(s string) (migrators.MigrationType, error) {
	if slices.Contains(AllMigrationTypes, migrators.MigrationType(strings.ToLower(s))) {
		return migrators.MigrationType(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unsupported migration type: %s", s)
}

const (
	ImportStatementTypeStatement migrators.ImportStatementType = "statement"
	ImportStatementTypeBlock     migrators.ImportStatementType = "block"
)

var AllImportStatementTypes = []migrators.ImportStatementType{
	ImportStatementTypeStatement,
	ImportStatementTypeBlock,
}

func ToImportStatementType(s string) (migrators.ImportStatementType, error) {
	if slices.Contains(AllImportStatementTypes, migrators.ImportStatementType(strings.ToLower(s))) {
		return migrators.ImportStatementType(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("invalid import statement type: %s", s)
}

type ExitCode int

const (
	ExitCodeSuccess ExitCode = iota
	ExitCodeFailedInputArgumentParsing
	ExitCodeFailedApiConnection
	ExitCodeFailedGeneratingTerraformOutput
)

type CLI struct {
	Args           []string
	StdOut, StdErr io.Writer
	StdIn          io.Reader
	Config         *migrators.Config
}

func NewDefaultCLI() *CLI {
	return &CLI{
		Args:   os.Args,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		StdIn:  os.Stdin,
	}
}

func (cli *CLI) Run() ExitCode {
	config, err := cli.parseInputArguments()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintf(cli.StdErr, "Error parsing input arguments: %v, run -h to get more information on running the script\n", err)
		return ExitCodeFailedInputArgumentParsing
	}
	cli.Config = config

	output, err := cli.generateOutput()
	if err != nil {
		_, _ = fmt.Fprintf(cli.StdErr, "Error generating output: %v\n", err)
		return ExitCodeFailedGeneratingTerraformOutput
	}
	_, _ = fmt.Fprint(cli.StdOut, output)

	return ExitCodeSuccess
}

func (cli *CLI) parseInputArguments() (*migrators.Config, error) {
	commandLine := flag.NewFlagSet(cli.Args[0], flag.ContinueOnError)
	commandLine.SetOutput(cli.StdErr)
	commandLine.Usage = func() {
		_, _ = fmt.Fprint(cli.StdErr, `DESCRIPTION:
    Migration script helps migrate from deprecated Rootly resources to newer equivalents.
    It fetches existing resources from your Rootly organization via the API, converts
    them to new Terraform resource formats, and generates both resource configurations
    and import statements to bring them under Terraform management.

USAGE:
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master [FLAGS] <migration_type>

MIGRATION TYPES:
    alert_routing_rules_to_alert_routes
        Migrates deprecated rootly_alert_routing_rule resources to rootly_alert_route
        resources. Fetches all alert routes from your Rootly organization and generates
        Terraform configurations with import statements.

FLAGS:
    -api-host string
        Rootly API host URL (default: https://api.rootly.com or ROOTLY_API_URL env var)
    -api-token string
        Rootly API authentication token (required, use ROOTLY_API_TOKEN env var)
    -import string
        Import statement format: "statement" or "block" (default: "statement")
        - "statement": generates terraform import shell commands
        - "block": generates import {} blocks for terraform plan/apply workflow

ENVIRONMENT VARIABLES:
    ROOTLY_API_TOKEN    Your Rootly API token (required)
    ROOTLY_API_URL      Custom Rootly API URL (optional, defaults to https://api.rootly.com)

EXAMPLES:
    # Basic usage with import commands
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master alert_routing_rules_to_alert_routes > alert_routes.tf

    # Generate import blocks instead of shell commands
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master -import=block alert_routing_rules_to_alert_routes > alert_routes.tf

    # With explicit API configuration
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master -api-host=https://api.rootly.com -api-token=your-token alert_routing_rules_to_alert_routes > alert_routes.tf

    # Set environment variables first
    export ROOTLY_API_TOKEN="your-api-token"
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@master alert_routing_rules_to_alert_routes > alert_routes.tf

WORKFLOW:
    1. Run the migration script to generate Terraform configuration
    2. Review the generated .tf file for accuracy
    3. Follow the instructions in the generated file to import resources
    4. Remove deprecated resources from your existing Terraform configuration
    5. Clean up old resources from Terraform state using 'terraform state rm'
`)
	}

	// flags
	importFlagString := commandLine.String("import", "statement", "Determines the output format for import statements")
	apiHost := commandLine.String("api-host", getEnvOrDefault("ROOTLY_API_URL", "https://api.rootly.com"), "Rootly API host")
	apiToken := commandLine.String("api-token", os.Getenv("ROOTLY_API_TOKEN"), "Rootly API token")

	if err := commandLine.Parse(cli.Args[1:]); err != nil {
		return nil, err
	}

	args := commandLine.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("no migration type specified, use -h for help")
	}

	// Parse migration type from positional arguments
	parsedMigrationType, err := ToMigrationType(args[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing migration type: %w", err)
	}

	importFlagType, err := ToImportStatementType(*importFlagString)
	if err != nil {
		return nil, fmt.Errorf("error parsing import flag: %w", err)
	}

	if *apiToken == "" {
		return nil, fmt.Errorf("api token is required, provide it via -api-token flag or ROOTLY_API_TOKEN environment variable")
	}

	return &migrators.Config{
		MigrationType: parsedMigrationType,
		ImportFlag:    importFlagType,
		ApiHost:       *apiHost,
		ApiToken:      *apiToken,
	}, nil
}

func (cli *CLI) generateOutput() (string, error) {
	switch cli.Config.MigrationType {
	case MigrationTypeAlertRoutingRulesToAlertRoutes:
		return migrators.HandleAlertRoutingRulesToAlertRoutes(cli.Config)
	default:
		return "", fmt.Errorf("unsupported migration type: %s, run -h to get more information on allowed migration types", cli.Config.MigrationType)
	}
}

func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}
