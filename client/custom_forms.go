package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type CustomForm struct {
	ID string `jsonapi:"primary,custom_forms"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Command string `jsonapi:"attr,command,omitempty"`
}

func (c *Client) ListCustomForms(params *rootlygo.ListCustomFormsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCustomFormsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	custom_forms, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CustomForm)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return custom_forms, nil
}

func (c *Client) CreateCustomForm(d *CustomForm) (*CustomForm, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_form: %s", err.Error())
	}

	req, err := rootlygo.NewCreateCustomFormRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom_form: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomForm))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_form: %s", err.Error())
	}

	return data.(*CustomForm), nil
}

func (c *Client) GetCustomForm(id string) (*CustomForm, error) {
	req, err := rootlygo.NewGetCustomFormRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get custom_form: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomForm))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_form: %s", err.Error())
	}

	return data.(*CustomForm), nil
}

func (c *Client) UpdateCustomForm(id string, custom_form *CustomForm) (*CustomForm, error) {
	buffer, err := MarshalData(custom_form)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom_form: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateCustomFormRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom_form: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomForm))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom_form: %s", err.Error())
	}

	return data.(*CustomForm), nil
}

func (c *Client) DeleteCustomForm(id string) error {
	req, err := rootlygo.NewDeleteCustomFormRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete custom_form: %s", err.Error())
	}

	return nil
}
