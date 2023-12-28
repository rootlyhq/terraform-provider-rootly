package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type IncidentPermissionSet struct {
	ID string `jsonapi:"primary,incident_permission_sets"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  PrivateIncidentPermissions []interface{} `jsonapi:"attr,private_incident_permissions,omitempty"`
  PublicIncidentPermissions []interface{} `jsonapi:"attr,public_incident_permissions,omitempty"`
}

func (c *Client) ListIncidentPermissionSets(params *rootlygo.ListIncidentPermissionSetsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentPermissionSetsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_permission_sets, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentPermissionSet)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_permission_sets, nil
}

func (c *Client) CreateIncidentPermissionSet(d *IncidentPermissionSet) (*IncidentPermissionSet, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentPermissionSetRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_permission_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSet))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set: %s", err.Error())
	}

	return data.(*IncidentPermissionSet), nil
}

func (c *Client) GetIncidentPermissionSet(id string) (*IncidentPermissionSet, error) {
	req, err := rootlygo.NewGetIncidentPermissionSetRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_permission_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSet))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set: %s", err.Error())
	}

	return data.(*IncidentPermissionSet), nil
}

func (c *Client) UpdateIncidentPermissionSet(id string, incident_permission_set *IncidentPermissionSet) (*IncidentPermissionSet, error) {
	buffer, err := MarshalData(incident_permission_set)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentPermissionSetRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_permission_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSet))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set: %s", err.Error())
	}

	return data.(*IncidentPermissionSet), nil
}

func (c *Client) DeleteIncidentPermissionSet(id string) error {
	req, err := rootlygo.NewDeleteIncidentPermissionSetRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_permission_set: %s", err.Error())
	}

	return nil
}
