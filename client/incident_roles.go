package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type IncidentRole struct {
	ID string `jsonapi:"primary,incident_roles"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Summary string `jsonapi:"attr,summary,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  Optional *bool `jsonapi:"attr,optional,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  AllowMultiUserAssignment *bool `jsonapi:"attr,allow_multi_user_assignment,omitempty"`
}

func (c *Client) ListIncidentRoles(params *rootlygo.ListIncidentRolesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentRolesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_roles, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentRole)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_roles, nil
}

func (c *Client) CreateIncidentRole(d *IncidentRole) (*IncidentRole, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_role: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentRoleRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role: %s", err.Error())
	}

	return data.(*IncidentRole), nil
}

func (c *Client) GetIncidentRole(id string) (*IncidentRole, error) {
	req, err := rootlygo.NewGetIncidentRoleRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role: %s", err.Error())
	}

	return data.(*IncidentRole), nil
}

func (c *Client) UpdateIncidentRole(id string, incident_role *IncidentRole) (*IncidentRole, error) {
	buffer, err := MarshalData(incident_role)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_role: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentRoleRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role: %s", err.Error())
	}

	return data.(*IncidentRole), nil
}

func (c *Client) DeleteIncidentRole(id string) error {
	req, err := rootlygo.NewDeleteIncidentRoleRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_role: %s", err.Error())
	}

	return nil
}
