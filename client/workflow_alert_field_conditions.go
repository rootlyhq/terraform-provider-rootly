package client

import (
	"encoding/json"
	"fmt"
	"io"

	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

// GetWorkflowAlertFieldConditions returns the alert field conditions of an alert workflow.
//
// The API accepts alert_field_conditions inside trigger_params on write, but stores them as a
// separate workflow_alert_field_conditions relationship that is only returned via
// ?include=alert_field_conditions.
func (c *Client) GetWorkflowAlertFieldConditions(id string) ([]interface{}, error) {
	include := rootlygo.GetWorkflowParamsIncludeAlertFieldConditions
	req, err := rootlygo.NewGetWorkflowRequest(c.Rootly.Server, id, &rootlygo.GetWorkflowParams{Include: &include})
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get workflow alert field conditions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var document struct {
		Included []struct {
			Type       string                 `json:"type"`
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"included"`
	}
	if err := json.Unmarshal(body, &document); err != nil {
		return nil, fmt.Errorf("Error decoding workflow alert field conditions: %w", err)
	}

	conditions := []interface{}{}
	for _, item := range document.Included {
		if item.Type != "workflow_alert_field_conditions" {
			continue
		}
		conditions = append(conditions, map[string]interface{}{
			"alert_field_id": item.Attributes["alert_field_id"],
			"condition_type": item.Attributes["condition_type"],
			"values":         item.Attributes["values"],
		})
	}

	return conditions, nil
}
