package client

import (
	"reflect"

	"fmt"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RetrospectiveConfiguration struct {
	ID              string        `jsonapi:"primary,retrospective_configurations"`
	Kind            string        `jsonapi:"attr,kind,omitempty"`
	SeverityIds     []interface{} `jsonapi:"attr,severity_ids,omitempty"`
	GroupIds        []interface{} `jsonapi:"attr,group_ids,omitempty"`
	IncidentTypeIds []interface{} `jsonapi:"attr,incident_type_ids,omitempty"`
}

func (c *Client) ListRetrospectiveConfigurations(params *rootlygo.ListRetrospectiveConfigurationsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveConfigurationsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	retrospective_configurations, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveConfiguration)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return retrospective_configurations, nil
}

func (c *Client) GetRetrospectiveConfiguration(id string) (*RetrospectiveConfiguration, error) {
	req, err := rootlygo.NewGetRetrospectiveConfigurationRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get retrospective_configuration: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveConfiguration))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling retrospective_configuration: %w", err)
	}

	return data.(*RetrospectiveConfiguration), nil
}

func (c *Client) UpdateRetrospectiveConfiguration(id string, retrospective_configuration *RetrospectiveConfiguration) (*RetrospectiveConfiguration, error) {
	buffer, err := MarshalData(retrospective_configuration)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling retrospective_configuration: %w", err)
	}

	req, err := rootlygo.NewUpdateRetrospectiveConfigurationRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update retrospective_configuration: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveConfiguration))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling retrospective_configuration: %w", err)
	}

	return data.(*RetrospectiveConfiguration), nil
}
