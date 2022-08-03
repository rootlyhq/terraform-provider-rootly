package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Environment struct {
	ID          string `jsonapi:"primary,environments"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Slug        string `jsonapi:"attr,slug,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
}

func (c *Client) CreateEnvironment(s *Environment) (*Environment, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling environment: %s", err.Error())
	}

	req, err := rootlygo.NewCreateEnvironmentRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create environment: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Environment))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling environment: %s", err.Error())
	}

	return data.(*Environment), nil
}

func (c *Client) GetEnvironment(id string) (*Environment, error) {
	req, err := rootlygo.NewGetEnvironmentRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get environment: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Environment))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling environment: %s", err.Error())
	}

	return data.(*Environment), nil
}

func (c *Client) UpdateEnvironment(id string, s *Environment) (*Environment, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling environment: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateEnvironmentRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update environment: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Environment))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling environment: %s", err.Error())
	}

	return data.(*Environment), nil
}

func (c *Client) DeleteEnvironment(id string) error {
	req, err := rootlygo.NewDeleteEnvironmentRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete environment: %s", id)
	}

	return nil
}
