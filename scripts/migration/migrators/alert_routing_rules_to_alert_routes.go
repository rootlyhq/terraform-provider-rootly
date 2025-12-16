package migrators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

type AlertRouteModel struct {
	Name            string
	Enabled         bool
	AlertsSourceIds []string
	OwningTeamIds   []string
	Rules           []AlertRouteRuleModel
}

type AlertRouteRuleModel struct {
	Name            string
	Destinations    []DestinationModel
	ConditionGroups []ConditionGroupModel
}

type DestinationModel struct {
	TargetType string
	TargetId   string
}

type ConditionGroupModel struct {
	Conditions []ConditionModel
}

type ConditionModel struct {
	PropertyFieldType          string
	PropertyFieldName          string
	PropertyFieldConditionType string
	PropertyFieldValue         string
	PropertyFieldValues        []string
	AlertUrgencyIds            []string
}

type RootlyClient struct {
	ApiHost  string
	ApiToken string
	client   *http.Client
}

func NewRootlyClient(apiHost, apiToken string) *RootlyClient {
	return &RootlyClient{
		ApiHost:  apiHost,
		ApiToken: apiToken,
		client:   &http.Client{},
	}
}

func (c *RootlyClient) FetchAlertRoutes() ([]client.AlertRoute, error) {
	var allRoutes []client.AlertRoute
	pageNumber := 1
	pageSize := 10

	for {
		url := fmt.Sprintf("%s/v1/alert_routes?page[number]=%d&page[size]=%d",
			strings.TrimSuffix(c.ApiHost, "/"), pageNumber, pageSize)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.ApiToken)
		req.Header.Set("Content-Type", "application/vnd.api+json")

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}

		var jsonApiResponse struct {
			Data []struct {
				ID         string `json:"id"`
				Attributes struct {
					Name            string        `json:"name"`
					Enabled         bool          `json:"enabled"`
					AlertsSourceIds []interface{} `json:"alerts_source_ids"`
					OwningTeamIds   []interface{} `json:"owning_team_ids"`
					Rules           []interface{} `json:"rules"`
				} `json:"attributes"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&jsonApiResponse); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("error decoding response: %w", err)
		}
		resp.Body.Close()

		if len(jsonApiResponse.Data) == 0 {
			break
		}

		for _, item := range jsonApiResponse.Data {
			route := client.AlertRoute{
				ID:              item.ID,
				Name:            item.Attributes.Name,
				Enabled:         &item.Attributes.Enabled,
				AlertsSourceIds: item.Attributes.AlertsSourceIds,
				OwningTeamIds:   item.Attributes.OwningTeamIds,
				Rules:           item.Attributes.Rules,
			}
			allRoutes = append(allRoutes, route)
		}

		pageNumber++
	}

	return allRoutes, nil
}

func ConvertAlertRouteToTerraform(route client.AlertRoute) AlertRouteModel {
	alertRoute := AlertRouteModel{
		Name:            route.Name,
		Enabled:         route.Enabled != nil && *route.Enabled,
		AlertsSourceIds: convertInterfaceArrayToStringArray(route.AlertsSourceIds),
		OwningTeamIds:   convertInterfaceArrayToStringArray(route.OwningTeamIds),
		Rules:           []AlertRouteRuleModel{},
	}

	if route.Rules != nil {
		for _, ruleInterface := range route.Rules {
			if ruleMap, ok := ruleInterface.(map[string]interface{}); ok {
				routeRule := AlertRouteRuleModel{
					Destinations:    []DestinationModel{},
					ConditionGroups: []ConditionGroupModel{},
				}

				if name, ok := ruleMap["name"].(string); ok {
					routeRule.Name = name
				}

				if destinations, ok := ruleMap["destinations"].([]interface{}); ok {
					for _, destInterface := range destinations {
						if destMap, ok := destInterface.(map[string]interface{}); ok {
							destination := DestinationModel{}
							if targetType, ok := destMap["target_type"].(string); ok {
								destination.TargetType = targetType
							}
							if targetId, ok := destMap["target_id"].(string); ok {
								destination.TargetId = targetId
							}
							routeRule.Destinations = append(routeRule.Destinations, destination)
						}
					}
				}

				if conditionGroups, ok := ruleMap["condition_groups"].([]interface{}); ok {
					for _, cgInterface := range conditionGroups {
						if cgMap, ok := cgInterface.(map[string]interface{}); ok {
							conditionGroup := ConditionGroupModel{
								Conditions: []ConditionModel{},
							}

							if conditions, ok := cgMap["conditions"].([]interface{}); ok {
								for _, condInterface := range conditions {
									if condMap, ok := condInterface.(map[string]interface{}); ok {
										condition := ConditionModel{}

										if val, ok := condMap["property_field_type"].(string); ok {
											condition.PropertyFieldType = val
										}
										if val, ok := condMap["property_field_name"].(string); ok {
											condition.PropertyFieldName = val
										}
										if val, ok := condMap["property_field_condition_type"].(string); ok {
											condition.PropertyFieldConditionType = val
										}
										if val, ok := condMap["property_field_value"].(string); ok {
											condition.PropertyFieldValue = val
										}

										if valArray, ok := condMap["property_field_values"].([]interface{}); ok {
											condition.PropertyFieldValues = convertInterfaceArrayToStringArray(valArray)
										}

										if urgencyArray, ok := condMap["alert_urgency_ids"].([]interface{}); ok {
											condition.AlertUrgencyIds = convertInterfaceArrayToStringArray(urgencyArray)
										}

										conditionGroup.Conditions = append(conditionGroup.Conditions, condition)
									}
								}
							}

							routeRule.ConditionGroups = append(routeRule.ConditionGroups, conditionGroup)
						}
					}
				}

				alertRoute.Rules = append(alertRoute.Rules, routeRule)
			}
		}
	}

	return alertRoute
}

func convertInterfaceArrayToStringArray(arr []interface{}) []string {
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func GenerateTerraformResource(resourceName string, alertRoute AlertRouteModel, ruleID string) (string, error) {
	tmplText := `resource "rootly_alert_route" "{{ .ResourceName }}" {
  name               = "{{ .AlertRoute.Name }}"
  enabled            = {{ .AlertRoute.Enabled }}
  alerts_source_ids  = [{{ range $i, $id := .AlertRoute.AlertsSourceIds }}{{ if $i }}, {{ end }}"{{ $id }}"{{ end }}]
{{ if .AlertRoute.OwningTeamIds }}  owning_team_ids    = [{{ range $i, $id := .AlertRoute.OwningTeamIds }}{{ if $i }}, {{ end }}"{{ $id }}"{{ end }}]
{{ end }}
{{ range .AlertRoute.Rules }}  rules {
    name = "{{ .Name }}"
{{ range .Destinations }}
    destinations {
      target_type = "{{ .TargetType }}"
      target_id   = "{{ .TargetId }}"
    }
{{ end }}
{{ range .ConditionGroups }}
    condition_groups {
{{ range .Conditions }}      conditions {
        property_field_type            = "{{ .PropertyFieldType }}"
        property_field_name            = "{{ .PropertyFieldName }}"
        property_field_condition_type  = "{{ .PropertyFieldConditionType }}"
{{ if .PropertyFieldValue }}        property_field_value           = "{{ .PropertyFieldValue }}"
{{ end }}{{ if .PropertyFieldValues }}        property_field_values          = [{{ range $i, $val := .PropertyFieldValues }}{{ if $i }}, {{ end }}"{{ $val }}"{{ end }}]
{{ end }}{{ if .AlertUrgencyIds }}        alert_urgency_ids              = [{{ range $i, $id := .AlertUrgencyIds }}{{ if $i }}, {{ end }}"{{ $id }}"{{ end }}]
{{ end }}      }
{{ end }}    }
{{ end }}  }
{{ end }}}`

	tmpl, err := template.New("terraform_resource").Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	data := struct {
		ResourceName string
		AlertRoute   AlertRouteModel
		RuleID       string
	}{
		ResourceName: resourceName,
		AlertRoute:   alertRoute,
		RuleID:       ruleID,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

func HandleAlertRoutingRulesToAlertRoutes(config *Config) (string, error) {
	client := NewRootlyClient(config.ApiHost, config.ApiToken)

	routes, err := client.FetchAlertRoutes()
	if err != nil {
		return "", fmt.Errorf("error fetching alert routes: %w", err)
	}

	var output strings.Builder
	var importStatements []string

	for i, route := range routes {
		alertRoute := ConvertAlertRouteToTerraform(route)
		resourceName := sanitizeResourceName(route.Name, i)
		resourceAddress := fmt.Sprintf("rootly_alert_route.%s", resourceName)

		resourceText, err := GenerateTerraformResource(resourceName, alertRoute, route.ID)
		if err != nil {
			return "", fmt.Errorf("error generating resource for route %s: %w", route.ID, err)
		}

		output.WriteString(resourceText)
		output.WriteString("\n")

		importStmt, err := GenerateImportStatement(config.ImportFlag, resourceAddress, route.ID)
		if err != nil {
			return "", fmt.Errorf("error generating import statement for route %s: %w", route.ID, err)
		}
		importStatements = append(importStatements, importStmt)
	}

	output.WriteString("\n# Import statements\n")
	for _, stmt := range importStatements {
		output.WriteString(stmt)
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString("# Instructions:\n")
	output.WriteString("# 1. Run 'terraform plan' to verify the import operations\n")
	output.WriteString("# 2. Run 'terraform apply' to apply the import operations\n")
	output.WriteString("# 3. Remove the import blocks/statements above from this file\n")
	output.WriteString("# 4. Remove deprecated 'rootly_alert_routing_rule' resources from your Terraform configuration\n")
	output.WriteString("# 5. Run 'terraform state rm <resource_address>' for each deprecated alert_routing_rule resource\n")
	output.WriteString("#    Example: terraform state rm rootly_alert_routing_rule.my_rule\n")

	return output.String(), nil
}

func sanitizeResourceName(name string, fallbackIndex int) string {
	if name == "" {
		return fmt.Sprintf("alert_route_%d", fallbackIndex)
	}

	result := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' {
			result += string(char)
		} else {
			result += "_"
		}
	}

	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		result = "_" + result
	}

	if result == "" {
		result = fmt.Sprintf("alert_route_%d", fallbackIndex)
	}

	return strings.ToLower(result)
}
