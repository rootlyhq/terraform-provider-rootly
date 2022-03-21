package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type Functionality struct {
	ID          string `jsonapi:"primary,functionalities"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Slug        string `jsonapi:"attr,slug,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) CreateFunctionality(f *Functionality) (*Functionality, error) {
	buffer, err := MarshalData(f)
	if err != nil {
		return nil, errors.Errorf("Error marshaling functionality: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFunctionalityRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create functionality: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) GetFunctionality(id string) (*Functionality, error) {
	req, err := rootlygo.NewGetFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get functionality: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) UpdateFunctionality(id string, f *Functionality) (*Functionality, error) {
	buffer, err := MarshalData(f)
	if err != nil {
		return nil, errors.Errorf("Error marshaling functionality: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFunctionalityRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update functionality: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) DeleteFunctionality(id string) error {
	req, err := rootlygo.NewDeleteFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete functionality: %s", id)
	}

	return nil
}
