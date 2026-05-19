package client

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/jsonapi"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type CustomFieldOption struct {
	ID            string `jsonapi:"primary,custom_field_options"`
	CustomFieldId int    `jsonapi:"attr,custom_field_id,omitempty"`
	Value         string `jsonapi:"attr,value,omitempty"`
	Color         string `jsonapi:"attr,color,omitempty"`
	Default       *bool  `jsonapi:"attr,default,omitempty"`
	Position      int    `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListCustomFieldOptions(id rootlygo_.ID, params *rootlygo.ListCustomFieldOptionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCustomFieldOptionsRequest(c.Rootly.Server, id.String(), params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	custom_field_options, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomFieldOption)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return custom_field_options, nil
}

func (c *Client) CreateCustomFieldOption(d *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling custom_field_option: %w", err)
	}

	req, err := rootlygo.NewCreateCustomFieldOptionRequestWithBody(c.Rootly.Server, strconv.Itoa(d.CustomFieldId), c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create custom_field_option: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field_option: %w", err)
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) GetCustomFieldOption(id rootlygo_.ID) (*CustomFieldOption, error) {
	req, err := rootlygo.NewGetCustomFieldOptionRequest(c.Rootly.Server, id.String())
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get custom_field_option: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field_option: %w", err)
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) UpdateCustomFieldOption(id rootlygo_.ID, custom_field_option *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(custom_field_option)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling custom_field_option: %w", err)
	}

	req, err := rootlygo.NewUpdateCustomFieldOptionRequestWithBody(c.Rootly.Server, id.String(), c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update custom_field_option: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling custom_field_option: %w", err)
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) DeleteCustomFieldOption(id rootlygo_.ID) error {
	req, err := rootlygo.NewDeleteCustomFieldOptionRequest(c.Rootly.Server, id.String())
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete custom_field_option: %w", err)
	}

	return nil
}
