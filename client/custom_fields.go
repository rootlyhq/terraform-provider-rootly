package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type CustomField struct {
	ID          string `jsonapi:"primary,custom_fields"`
	Label       string `jsonapi:"attr,label,omitempty"`
	Kind        string `jsonapi:"attr,kind,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Enabled     *bool  `jsonapi:"attr,enabled,omitempty"`
	Shown       []interface{} `jsonapi:"attr,shown,omitempty"`
	Required    []interface{} `jsonapi:"attr,required,omitempty"`
}

func (c *Client) CreateCustomField(i *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom field: %s", err.Error())
	}

	req, err := rootlygo.NewCreateCustomFieldRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field: %s", err.Error())
	}

	return data.(*CustomField), nil
}

func (c *Client) GetCustomField(id string) (*CustomField, error) {
	req, err := rootlygo.NewGetCustomFieldRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get custom field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field: %s", err.Error())
	}

	return data.(*CustomField), nil
}

func (c *Client) UpdateCustomField(id string, i *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom field: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateCustomFieldRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field: %s", err.Error())
	}

	return data.(*CustomField), nil
}

func (c *Client) DeleteCustomField(id string) error {
	req, err := rootlygo.NewDeleteCustomFieldRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete custom field: %s", id)
	}

	return nil
}
