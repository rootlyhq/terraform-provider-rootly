package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type CustomField struct {
	ID string `jsonapi:"primary,custom_fields"`
	Label string `jsonapi:"attr,label,omitempty"`
  Kind string `jsonapi:"attr,kind,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Shown []interface{} `jsonapi:"attr,shown,omitempty"`
  Required []interface{} `jsonapi:"attr,required,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListCustomFields(params *rootlygo.ListCustomFieldsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCustomFieldsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	custom_fields, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomField)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return custom_fields, nil
}

func (c *Client) CreateCustomField(d *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_field: %s", err.Error())
	}

	req, err := rootlygo.NewCreateCustomFieldRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom_field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field: %s", err.Error())
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
		return nil, errors.Errorf("Failed to make request to get custom_field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field: %s", err.Error())
	}

	return data.(*CustomField), nil
}

func (c *Client) UpdateCustomField(id string, custom_field *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(custom_field)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_field: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateCustomFieldRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom_field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field: %s", err.Error())
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
		return errors.Errorf("Failed to make request to delete custom_field: %s", id)
	}

	return nil
}
