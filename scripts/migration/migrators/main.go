package migrators

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type MigrationType string
type ImportStatementType string

const (
	ImportStatementTypeStatement ImportStatementType = "statement"
	ImportStatementTypeBlock     ImportStatementType = "block"
)

type Config struct {
	MigrationType MigrationType
	ImportFlag    ImportStatementType
	ApiHost       string
	ApiToken      string
}

// ImportModel represents the data for terraform import statements
type ImportModel struct {
	ResourceAddress string
	Id              string
	IdEscaped       string
}

func GenerateImportStatement(importType ImportStatementType, resourceAddress, resourceID string) (string, error) {
	var templateText string

	switch importType {
	case ImportStatementTypeStatement:
		templateText = `# terraform import {{ .ResourceAddress }} '{{ .Id }}'`
	case ImportStatementTypeBlock:
		templateText = `import {
  to = {{ .ResourceAddress }}
  id = "{{ .IdEscaped }}"
}`
	default:
		return "", fmt.Errorf("unsupported import type: %s", importType)
	}

	tmpl, err := template.New("import").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("error parsing import template: %w", err)
	}

	// Escape the ID for block format if needed
	idEscaped := resourceID
	if importType == ImportStatementTypeBlock {
		idEscaped = strings.ReplaceAll(resourceID, `"`, `\"`)
	}

	data := ImportModel{
		ResourceAddress: resourceAddress,
		Id:              resourceID,
		IdEscaped:       idEscaped,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing import template: %w", err)
	}

	return buf.String(), nil
}

