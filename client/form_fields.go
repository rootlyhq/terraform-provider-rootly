package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type FormField struct {
	ID string `jsonapi:"primary,form_fields"`
	Kind string `jsonapi:"attr,kind,omitempty"`
  InputKind string `jsonapi:"attr,input_kind,omitempty"`
  ValueKind string `jsonapi:"attr,value_kind,omitempty"`
  Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Shown []interface{} `jsonapi:"attr,shown,omitempty"`
  Required []interface{} `jsonapi:"attr,required,omitempty"`
  ShowOnIncidentDetails *bool `jsonapi:"attr,show_on_incident_details,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  DefaultValues []interface{} `jsonapi:"attr,default_values,omitempty"`
}

func (c *Client) ListFormFields(params *rootlygo.ListFormFieldsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormFieldsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_fields, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormField)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_fields, nil
}

func (c *Client) CreateFormField(d *FormField) (*FormField, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormFieldRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field: %s", err.Error())
	}

	return data.(*FormField), nil
}

func (c *Client) GetFormField(id string) (*FormField, error) {
	req, err := rootlygo.NewGetFormFieldRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field: %s", err.Error())
	}

	return data.(*FormField), nil
}

func (c *Client) UpdateFormField(id string, form_field *FormField) (*FormField, error) {
	buffer, err := MarshalData(form_field)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormFieldRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field: %s", err.Error())
	}

	return data.(*FormField), nil
}

func (c *Client) DeleteFormField(id string) error {
	req, err := rootlygo.NewDeleteFormFieldRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_field: %s", err.Error())
	}

	return nil
}
