package client

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type CustomField struct {
	ID          string        `jsonapi:"primary,custom_fields"`
	Label       string        `jsonapi:"attr,label,omitempty"`
	Kind        string        `jsonapi:"attr,kind,omitempty"`
	Enabled     *bool         `jsonapi:"attr,enabled,omitempty"`
	Slug        string        `jsonapi:"attr,slug,omitempty"`
	Description string        `jsonapi:"attr,description,omitempty"`
	Shown       []interface{} `jsonapi:"attr,shown,omitempty"`
	Required    []interface{} `jsonapi:"attr,required,omitempty"`
	Default     string        `jsonapi:"attr,default,omitempty"`
	Position    int           `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListCustomFields(params *rootlygo.ListCustomFieldsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCustomFieldsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	custom_fields, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomField)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return custom_fields, nil
}

func (c *Client) CreateCustomField(d *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling custom_field: %w", err)
	}

	req, err := rootlygo.NewCreateCustomFieldRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create custom_field: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field: %w", err)
	}

	return data.(*CustomField), nil
}

func (c *Client) GetCustomField(id rootlygo_.ID) (*CustomField, error) {
	req, err := rootlygo.NewGetCustomFieldRequest(c.Rootly.Server, id.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get custom_field: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field: %w", err)
	}

	return data.(*CustomField), nil
}

func (c *Client) UpdateCustomField(id rootlygo_.ID, custom_field *CustomField) (*CustomField, error) {
	buffer, err := MarshalData(custom_field)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling custom_field: %w", err)
	}

	req, err := rootlygo.NewUpdateCustomFieldRequestWithBody(c.Rootly.Server, id.String(), c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update custom_field: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomField))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field: %w", err)
	}

	return data.(*CustomField), nil
}

func (c *Client) DeleteCustomField(id rootlygo_.ID) error {
	req, err := rootlygo.NewDeleteCustomFieldRequest(c.Rootly.Server, id.String())
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete custom_field: %w", err)
	}

	return nil
}
