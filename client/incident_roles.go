package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type IncidentRole struct {
	ID          string `jsonapi:"primary,incident_roles"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Summary     string `jsonapi:"attr,summary,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Enabled     *bool  `jsonapi:"attr,enabled,omitempty"`
}

func (c *Client) CreateIncidentRole(i *IncidentRole) (*IncidentRole, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident role: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentRoleRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident role: %s", err.Error())
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
		return nil, errors.Errorf("Failed to make request to get incident role: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident role: %s", err.Error())
	}

	return data.(*IncidentRole), nil
}

func (c *Client) UpdateIncidentRole(id string, i *IncidentRole) (*IncidentRole, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident role: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentRoleRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident role: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRole))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident role: %s", err.Error())
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
		return errors.Errorf("Failed to make request to delete incident role: %s", id)
	}

	return nil
}
