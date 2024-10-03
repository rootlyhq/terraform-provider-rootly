package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type FormSetCondition struct {
	ID string `jsonapi:"primary,form_set_conditions"`
	FormSetId string `jsonapi:"attr,form_set_id,omitempty"`
  FormFieldId string `jsonapi:"attr,form_field_id,omitempty"`
  Comparison string `jsonapi:"attr,comparison,omitempty"`
  Values []interface{} `jsonapi:"attr,values,omitempty"`
}

func (c *Client) ListFormSetConditions(id string, params *rootlygo.ListFormSetConditionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormSetConditionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_set_conditions, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormSetCondition)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_set_conditions, nil
}

func (c *Client) CreateFormSetCondition(d *FormSetCondition) (*FormSetCondition, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_set_condition: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormSetConditionRequestWithBody(c.Rootly.Server, d.FormSetId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_set_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSetCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set_condition: %s", err.Error())
	}

	return data.(*FormSetCondition), nil
}

func (c *Client) GetFormSetCondition(id string) (*FormSetCondition, error) {
	req, err := rootlygo.NewGetFormSetConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_set_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSetCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set_condition: %s", err.Error())
	}

	return data.(*FormSetCondition), nil
}

func (c *Client) UpdateFormSetCondition(id string, form_set_condition *FormSetCondition) (*FormSetCondition, error) {
	buffer, err := MarshalData(form_set_condition)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_set_condition: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormSetConditionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_set_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormSetCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_set_condition: %s", err.Error())
	}

	return data.(*FormSetCondition), nil
}

func (c *Client) DeleteFormSetCondition(id string) error {
	req, err := rootlygo.NewDeleteFormSetConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_set_condition: %s", err.Error())
	}

	return nil
}
