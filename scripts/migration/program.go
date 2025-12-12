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

type ObjectType string

const (
	ObjectTypeAlertRoutes ObjectType = "alert_routes"
)

var AllObjectTypes = []ObjectType{
	ObjectTypeAlertRoutes,
}

func ToObjectType(s string) (ObjectType, error) {
	if slices.Contains(AllObjectTypes, ObjectType(strings.ToLower(s))) {
		return ObjectType(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unsupported object type: %s", s)
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
	ObjectType ObjectType
	ImportFlag ImportStatementType
	ApiHost    string
	ApiToken   string
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
		_, _ = fmt.Fprint(p.StdErr, `Migration script's purpose is to generate terraform resources from existing Rootly objects.
It fetches alert routes from the Rootly API and converts them to terraform alert route resources.
The script writes generated terraform resources to STDOUT in case you want to redirect it to a file.
Any logs or errors are written to STDERR.

usage: migration_script [-import=<statement|block>] [-api-host=<api_host>] [-api-token=<api_token>] <object_type>

api-host optional flag specifies the Rootly API host (defaults to https://api.rootly.com or ROOTLY_API_URL env var)
api-token optional flag specifies the Rootly API token (defaults to ROOTLY_API_TOKEN env var)
import optional flag determines the output format for import statements. The possible values are:
	- "statement" will print appropriate terraform import command at the end of generated content (default)
	- "block" will generate import block at the end of generated content

object_type represents the type of Rootly object you want to migrate.
	Currently supported object types are:
		- "alert_routes" which fetches all alert routes and converts them to terraform resources

example usage:
	migration_script alert_routes > ./alert_routes.tf
	migration_script -import=block -api-token=your-token alert_routes > ./alert_routes.tf
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
		return nil, fmt.Errorf("no object type specified, use -h for help")
	}
	parsedObjectType, err := ToObjectType(args[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing object type: %w", err)
	}

	importFlagType, err := ToImportStatementType(*importFlagString)
	if err != nil {
		return nil, fmt.Errorf("error parsing import flag: %w", err)
	}

	if *apiToken == "" {
		return nil, fmt.Errorf("api token is required, provide it via -api-token flag or ROOTLY_API_TOKEN environment variable")
	}

	return &Config{
		ObjectType: parsedObjectType,
		ImportFlag: importFlagType,
		ApiHost:    *apiHost,
		ApiToken:   *apiToken,
	}, nil
}

func (p *Program) generateOutput() (string, error) {
	switch p.Config.ObjectType {
	case ObjectTypeAlertRoutes:
		return HandleAlertRoutes(p.Config)
	default:
		return "", fmt.Errorf("unsupported object type: %s, run -h to get more information on allowed object types", p.Config.ObjectType)
	}
}

func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}