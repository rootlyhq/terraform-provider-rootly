package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RetrospectiveConfiguration struct {
	ID string `jsonapi:"primary,retrospective_configurations"`
	Kind string `jsonapi:"attr,kind,omitempty"`
	SeverityIds []interface{} `jsonapi:"attr,severity_ids,omitempty"`
	GroupIds []interface{} `jsonapi:"attr,group_ids,omitempty"`
	IncidentTypeIds []interface{} `jsonapi:"attr,incident_type_ids,omitempty"`
}

func (c *Client) ListRetrospectiveConfigurations(params *rootlygo.ListRetrospectiveConfigurationsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveConfigurationsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	retrospective_configurations, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveConfiguration)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return retrospective_configurations, nil
}

func (c *Client) GetRetrospectiveConfiguration(id string) (*RetrospectiveConfiguration, error) {
	req, err := rootlygo.NewGetRetrospectiveConfigurationRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get retrospective_configuration: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveConfiguration))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_configuration: %s", err.Error())
	}

	return data.(*RetrospectiveConfiguration), nil
}

func (c *Client) UpdateRetrospectiveConfiguration(id string, retrospective_configuration *RetrospectiveConfiguration) (*RetrospectiveConfiguration, error) {
	buffer, err := MarshalData(retrospective_configuration)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_configuration: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRetrospectiveConfigurationRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update retrospective_configuration: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveConfiguration))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_configuration: %s", err.Error())
	}

	return data.(*RetrospectiveConfiguration), nil
}
