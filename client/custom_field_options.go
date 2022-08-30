package client

import (
	"reflect"
	"strconv"
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type CustomFieldOption struct {
	ID string `jsonapi:"primary,custom_field_options"`
	CustomFieldId int `jsonapi:"attr,custom_field_id,omitempty"`
  Value string `jsonapi:"attr,value,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListCustomFieldOptions(id string, params *rootlygo.ListCustomFieldOptionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCustomFieldOptionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	custom_field_options, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomFieldOption)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return custom_field_options, nil
}

func (c *Client) CreateCustomFieldOption(d *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_field_option: %s", err.Error())
	}

	req, err := rootlygo.NewCreateCustomFieldOptionRequestWithBody(c.Rootly.Server, strconv.Itoa(d.CustomFieldId), c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom_field_option: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field_option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) GetCustomFieldOption(id string) (*CustomFieldOption, error) {
	req, err := rootlygo.NewGetCustomFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get custom_field_option: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field_option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) UpdateCustomFieldOption(id string, custom_field_option *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(custom_field_option)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_field_option: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateCustomFieldOptionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom_field_option: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_field_option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) DeleteCustomFieldOption(id string) error {
	req, err := rootlygo.NewDeleteCustomFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete custom_field_option: %s", id)
	}

	return nil
}
