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

type Program struct {
	Args           []string
	StdOut, StdErr io.Writer
	StdIn          io.Reader
	Config         *migrators.Config
}

func NewDefaultProgram() *Program {
	return &Program{
		Args:   os.Args,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		StdIn:  os.Stdin,
	}
}

func (p *Program) Run() ExitCode {
	config, err := p.parseInputArguments()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintf(p.StdErr, "Error parsing input arguments: %v, run -h to get more information on running the script\n", err)
		return ExitCodeFailedInputArgumentParsing
	}
	p.Config = config

	output, err := p.generateOutput()
	if err != nil {
		_, _ = fmt.Fprintf(p.StdErr, "Error generating output: %v\n", err)
		return ExitCodeFailedGeneratingTerraformOutput
	}
	_, _ = fmt.Fprint(p.StdOut, output)

	return ExitCodeSuccess
}

func (p *Program) parseInputArguments() (*migrators.Config, error) {
	commandLine := flag.NewFlagSet(p.Args[0], flag.ContinueOnError)
	commandLine.SetOutput(p.StdErr)
	commandLine.Usage = func() {
		_, _ = fmt.Fprint(p.StdErr, `Migration script's purpose is to help migrate deprecated Rootly resources to newer equivalents.
It fetches existing resources from the Rootly API and converts them to new resource formats.
The script writes generated terraform resources to STDOUT which is usually redirected to a file.

USAGE:
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main [FLAGS] <migration_type>

FLAGS:
    -api-host string
        Rootly API host (defaults to https://api.rootly.com or ROOTLY_API_URL env var)
    -api-token string
        Rootly API token (defaults to ROOTLY_API_TOKEN env var)
    -import string
        Output format for import statements: "statement" or "block" (default "statement")
        - "statement": prints terraform import commands
        - "block": generates import blocks for use with Terraform

MIGRATION TYPES:
    alert_routing_rules_to_alert_routes
        Migrates deprecated alert_routing_rule resources to alert_route resources

EXAMPLES:
    # Basic usage
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main alert_routing_rules_to_alert_routes > ./alert_routes.tf

    # With import blocks
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -import=block alert_routing_rules_to_alert_routes > ./alert_routes.tf

    # With custom API configuration
    go run github.com/rootlyhq/terraform-provider-rootly/v2/scripts/migration@main -api-host=https://api.rootly.com -api-token=your-token alert_routing_rules_to_alert_routes > ./alert_routes.tf

ENVIRONMENT VARIABLES:
    ROOTLY_API_TOKEN    Your Rootly API token
    ROOTLY_API_URL      Custom Rootly API URL (optional, defaults to https://api.rootly.com)
`)
	}

	// flags
	importFlagString := commandLine.String("import", "statement", "Determines the output format for import statements")
	apiHost := commandLine.String("api-host", getEnvOrDefault("ROOTLY_API_URL", "https://api.rootly.com"), "Rootly API host")
	apiToken := commandLine.String("api-token", os.Getenv("ROOTLY_API_TOKEN"), "Rootly API token")

	if err := commandLine.Parse(p.Args[1:]); err != nil {
		return nil, err
	}

	// positional arguments
	args := commandLine.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("no migration type specified, use -h for help")
	}
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

func (p *Program) generateOutput() (string, error) {
	switch p.Config.MigrationType {
	case MigrationTypeAlertRoutingRulesToAlertRoutes:
		return migrators.HandleAlertRoutingRulesToAlertRoutes(p.Config)
	default:
		return "", fmt.Errorf("unsupported migration type: %s, run -h to get more information on allowed migration types", p.Config.MigrationType)
	}
}

func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}
