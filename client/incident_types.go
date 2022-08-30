package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type IncidentType struct {
	ID string `jsonapi:"primary,incident_types"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) ListIncidentTypes(params *rootlygo.ListIncidentTypesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentTypesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_types, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentType)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_types, nil
}

func (c *Client) CreateIncidentType(d *IncidentType) (*IncidentType, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_type: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentTypeRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_type: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_type: %s", err.Error())
	}

	return data.(*IncidentType), nil
}

func (c *Client) GetIncidentType(id string) (*IncidentType, error) {
	req, err := rootlygo.NewGetIncidentTypeRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_type: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_type: %s", err.Error())
	}

	return data.(*IncidentType), nil
}

func (c *Client) UpdateIncidentType(id string, incident_type *IncidentType) (*IncidentType, error) {
	buffer, err := MarshalData(incident_type)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_type: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentTypeRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_type: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_type: %s", err.Error())
	}

	return data.(*IncidentType), nil
}

func (c *Client) DeleteIncidentType(id string) error {
	req, err := rootlygo.NewDeleteIncidentTypeRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_type: %s", id)
	}

	return nil
}
