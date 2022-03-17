package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type IncidentType struct {
	ID          string `jsonapi:"primary,incident_types"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) CreateIncidentType(i *IncidentType) (*IncidentType, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident type: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentTypeRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident type: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident type: %s", err.Error())
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
		return nil, errors.Errorf("Failed to make request to get incident type: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident type: %s", err.Error())
	}

	return data.(*IncidentType), nil
}

func (c *Client) UpdateIncidentType(id string, i *IncidentType) (*IncidentType, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident type: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentTypeRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident type: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentType))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident type: %s", err.Error())
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
		return errors.Errorf("Failed to make request to delete incident type: %s", id)
	}

	return nil
}
