// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package client

import (
    "fmt"
	"reflect"
	
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type FormField struct {
	ID string `jsonapi:"primary,form_fields"`
	Kind string `jsonapi:"attr,kind,omitempty"`
  InputKind string `jsonapi:"attr,input_kind,omitempty"`
  ValueKind string `jsonapi:"attr,value_kind,omitempty"`
  ValueKindCatalogId string `jsonapi:"attr,value_kind_catalog_id,omitempty"`
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
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	form_fields, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormField)))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return form_fields, nil
}

func (c *Client) CreateFormField(d *FormField) (*FormField, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling form_field: %w", err)
	}

	req, err := rootlygo.NewCreateFormFieldRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create form_field: %s", err)
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling form_field: %w", err)
	}

	return data.(*FormField), nil
}

func (c *Client) GetFormField(id string) (*FormField, error) {
	req, err := rootlygo.NewGetFormFieldRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get form_field: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling form_field: %w", err)
	}

	return data.(*FormField), nil
}

func (c *Client) UpdateFormField(id string, form_field *FormField) (*FormField, error) {
	buffer, err := MarshalData(form_field)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling form_field: %w", err)
	}

	req, err := rootlygo.NewUpdateFormFieldRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update form_field: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(FormField))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling form_field: %w", err)
	}

	return data.(*FormField), nil
}

func (c *Client) DeleteFormField(id string) error {
	req, err := rootlygo.NewDeleteFormFieldRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete form_field: %w", err)
	}

	return nil
}
