package client

import (
	"bytes"
	"github.com/google/jsonapi"
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

func (f Functionality) Marshal() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, f); err != nil {
		return nil, errors.Errorf("Error marshaling functionality (creation): %s", err.Error())
	}

	return buffer, nil
}

func (c *Client) CreateFunctionality(f *Functionality) (*Functionality, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, f); err != nil {
		return nil, errors.Errorf("Error marshaling functionality (creation): %s", err.Error())
	}

	req, err := rootlygo.NewCreateFunctionalityRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality (creation): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create functionality" + "\n\n" + buffer.String())
	}

	functionality := new(Functionality)
	if err := jsonapi.UnmarshalPayload(resp.Body, functionality); err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality (creation): %s", err.Error())
	}

	return functionality, nil
}

func (c *Client) GetFunctionality(id string) (*Functionality, error) {
	req, err := rootlygo.NewGetFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get functionality: %s", id)
	}

	functionality := new(Functionality)
	if err := jsonapi.UnmarshalPayload(resp.Body, functionality); err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality (read): %s", err.Error())
	}

	return functionality, nil
}

func (c *Client) UpdateFunctionality(id string, f *Functionality) (*Functionality, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, f); err != nil {
		return nil, errors.Errorf("Error marshaling functionality (update): %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFunctionalityRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality (update): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update functionality: %s", id)
	}

	functionality := new(Functionality)
	if err := jsonapi.UnmarshalPayload(resp.Body, functionality); err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality (update): %s", err.Error())
	}

	return functionality, nil
}

func (c *Client) DeleteFunctionality(id string) error {
	req, err := rootlygo.NewDeleteFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error unmarshaling functionality (delete): %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete functionality: %s", id)
	}

	return nil
}
