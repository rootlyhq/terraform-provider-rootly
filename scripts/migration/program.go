package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type MigrationType string

const (
	MigrationTypeAlertRoutingRulesToAlertRoutes MigrationType = "alert_routing_rules_to_alert_routes"
)

var AllMigrationTypes = []MigrationType{
	MigrationTypeAlertRoutingRulesToAlertRoutes,
}

func ToMigrationType(s string) (MigrationType, error) {
	if slices.Contains(AllMigrationTypes, MigrationType(strings.ToLower(s))) {
		return MigrationType(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unsupported migration type: %s", s)
}

type ImportStatementType string

const (
	ImportStatementTypeStatement ImportStatementType = "statement"
	ImportStatementTypeBlock     ImportStatementType = "block"
)

var AllImportStatementTypes = []ImportStatementType{
	ImportStatementTypeStatement,
	ImportStatementTypeBlock,
}

func ToImportStatementType(s string) (ImportStatementType, error) {
	if slices.Contains(AllImportStatementTypes, ImportStatementType(strings.ToLower(s))) {
		return ImportStatementType(strings.ToLower(s)), nil
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

type Config struct {
	MigrationType MigrationType
	ImportFlag    ImportStatementType
	ApiHost       string
	ApiToken      string
}

type Program struct {
	Args           []string
	StdOut, StdErr io.Writer
	StdIn          io.Reader
	Config         *Config
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

func (p *Program) parseInputArguments() (*Config, error) {
	commandLine := flag.NewFlagSet(p.Args[0], flag.ContinueOnError)
	commandLine.SetOutput(p.StdErr)
	commandLine.Usage = func() {
		_, _ = fmt.Fprint(p.StdErr, `Migration script's purpose is to help migrate deprecated Rootly resources to newer equivalents.
It fetches existing resources from the Rootly API and converts them to new resource formats.
The script writes generated terraform resources to STDOUT in case you want to redirect it to a file.
Any logs or errors are written to STDERR.

usage: ` + p.Args[0] + ` [-import=<statement|block>] [-api-host=<api_host>] [-api-token=<api_token>] <migration_type>

api-host optional flag specifies the Rootly API host (defaults to https://api.rootly.com or ROOTLY_API_URL env var)
api-token optional flag specifies the Rootly API token (defaults to ROOTLY_API_TOKEN env var)
import optional flag determines the output format for import statements. The possible values are:
	- "statement" will print appropriate terraform import command at the end of generated content (default)
	- "block" will generate import block at the end of generated content

migration_type represents the type of migration to perform. Currently supported migrations are:
	- "alert_routing_rules_to_alert_routes" migrates deprecated alert_routing_rule resources to alert_route resources

example usage:
	` + p.Args[0] + ` alert_routing_rules_to_alert_routes > ./alert_routes.tf
	` + p.Args[0] + ` -import=block -api-token=your-token alert_routing_rules_to_alert_routes > ./alert_routes.tf
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

	return &Config{
		MigrationType: parsedMigrationType,
		ImportFlag:    importFlagType,
		ApiHost:       *apiHost,
		ApiToken:      *apiToken,
	}, nil
}

func (p *Program) generateOutput() (string, error) {
	switch p.Config.MigrationType {
	case MigrationTypeAlertRoutingRulesToAlertRoutes:
		return HandleAlertRoutingRulesToAlertRoutes(p.Config)
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