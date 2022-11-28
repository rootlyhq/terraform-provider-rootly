package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type FormFieldOption struct {
	ID string `jsonapi:"primary,form_field_options"`
	FormFieldId string `jsonapi:"attr,form_field_id,omitempty"`
  Value string `jsonapi:"attr,value,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
  Default *bool `jsonapi:"attr,default,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListFormFieldOptions(id string, params *rootlygo.ListFormFieldOptionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormFieldOptionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_field_options, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormFieldOption)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_field_options, nil
}

func (c *Client) CreateFormFieldOption(d *FormFieldOption) (*FormFieldOption, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_option: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormFieldOptionRequestWithBody(c.Rootly.Server, d.FormFieldId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_field_option: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_option: %s", err.Error())
	}

	return data.(*FormFieldOption), nil
}

func (c *Client) GetFormFieldOption(id string) (*FormFieldOption, error) {
	req, err := rootlygo.NewGetFormFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_field_option: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_option: %s", err.Error())
	}

	return data.(*FormFieldOption), nil
}

func (c *Client) UpdateFormFieldOption(id string, form_field_option *FormFieldOption) (*FormFieldOption, error) {
	buffer, err := MarshalData(form_field_option)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_option: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormFieldOptionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_field_option: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_option: %s", err.Error())
	}

	return data.(*FormFieldOption), nil
}

func (c *Client) DeleteFormFieldOption(id string) error {
	req, err := rootlygo.NewDeleteFormFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_field_option: %s", err.Error())
	}

	return nil
}
