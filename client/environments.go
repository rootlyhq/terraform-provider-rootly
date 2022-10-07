package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Environment struct {
	ID string `jsonapi:"primary,environments"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) ListEnvironments(params *rootlygo.ListEnvironmentsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListEnvironmentsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	environments, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Environment)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return environments, nil
}

func (c *Client) CreateEnvironment(d *Environment) (*Environment, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling environment: %s", err.Error())
	}

	req, err := rootlygo.NewCreateEnvironmentRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create environment: %s", err.Error())
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
		return nil, errors.Errorf("Failed to make request to get environment: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Environment))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling environment: %s", err.Error())
	}

	return data.(*Environment), nil
}

func (c *Client) UpdateEnvironment(id string, environment *Environment) (*Environment, error) {
	buffer, err := MarshalData(environment)
	if err != nil {
		return nil, errors.Errorf("Error marshaling environment: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateEnvironmentRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update environment: %s", err.Error())
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
		return errors.Errorf("Failed to make request to delete environment: %s", err.Error())
	}

	return nil
}
