package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type IncidentPermissionSetBoolean struct {
	ID string `jsonapi:"primary,incident_permission_set_booleans"`
	IncidentPermissionSetId string `jsonapi:"attr,incident_permission_set_id,omitempty"`
  Kind string `jsonapi:"attr,kind,omitempty"`
  Private *bool `jsonapi:"attr,private,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
}

func (c *Client) ListIncidentPermissionSetBooleans(id string, params *rootlygo.ListIncidentPermissionSetBooleansParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentPermissionSetBooleansRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_permission_set_booleans, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentPermissionSetBoolean)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_permission_set_booleans, nil
}

func (c *Client) CreateIncidentPermissionSetBoolean(d *IncidentPermissionSetBoolean) (*IncidentPermissionSetBoolean, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set_boolean: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentPermissionSetBooleanRequestWithBody(c.Rootly.Server, d.IncidentPermissionSetId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_permission_set_boolean: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetBoolean))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_boolean: %s", err.Error())
	}

	return data.(*IncidentPermissionSetBoolean), nil
}

func (c *Client) GetIncidentPermissionSetBoolean(id string) (*IncidentPermissionSetBoolean, error) {
	req, err := rootlygo.NewGetIncidentPermissionSetBooleanRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_permission_set_boolean: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetBoolean))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_boolean: %s", err.Error())
	}

	return data.(*IncidentPermissionSetBoolean), nil
}

func (c *Client) UpdateIncidentPermissionSetBoolean(id string, incident_permission_set_boolean *IncidentPermissionSetBoolean) (*IncidentPermissionSetBoolean, error) {
	buffer, err := MarshalData(incident_permission_set_boolean)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set_boolean: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentPermissionSetBooleanRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_permission_set_boolean: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetBoolean))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_boolean: %s", err.Error())
	}

	return data.(*IncidentPermissionSetBoolean), nil
}

func (c *Client) DeleteIncidentPermissionSetBoolean(id string) error {
	req, err := rootlygo.NewDeleteIncidentPermissionSetBooleanRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_permission_set_boolean: %s", err.Error())
	}

	return nil
}
