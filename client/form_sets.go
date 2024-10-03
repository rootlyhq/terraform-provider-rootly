package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type FormSet struct {
	ID string `jsonapi:"primary,form_sets"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  IsDefault *bool `jsonapi:"attr,is_default,omitempty"`
  Forms []interface{} `jsonapi:"attr,forms,omitempty"`
}

func (c *Client) ListFormSets(params *rootlygo.ListFormSetsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormSetsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_sets, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormSet)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_sets, nil
}

func (c *Client) CreateFormSet(d *FormSet) (*FormSet, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_set: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormSetRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSet))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set: %s", err.Error())
	}

	return data.(*FormSet), nil
}

func (c *Client) GetFormSet(id string) (*FormSet, error) {
	req, err := rootlygo.NewGetFormSetRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSet))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set: %s", err.Error())
	}

	return data.(*FormSet), nil
}

func (c *Client) UpdateFormSet(id string, form_set *FormSet) (*FormSet, error) {
	buffer, err := MarshalData(form_set)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_set: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormSetRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_set: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSet))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set: %s", err.Error())
	}

	return data.(*FormSet), nil
}

func (c *Client) DeleteFormSet(id string) error {
	req, err := rootlygo.NewDeleteFormSetRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_set: %s", err.Error())
	}

	return nil
}
