// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type FormFieldPlacementCondition struct {
	ID string `jsonapi:"primary,form_field_placement_conditions"`
	FormFieldPlacementId string `jsonapi:"attr,form_field_placement_id,omitempty"`
  Conditioned string `jsonapi:"attr,conditioned,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  FormFieldId string `jsonapi:"attr,form_field_id,omitempty"`
  Comparison string `jsonapi:"attr,comparison,omitempty"`
  Values []interface{} `jsonapi:"attr,values,omitempty"`
}

func (c *Client) ListFormFieldPlacementConditions(id string, params *rootlygo.ListFormFieldPlacementConditionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormFieldPlacementConditionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_field_placement_conditions, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormFieldPlacementCondition)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_field_placement_conditions, nil
}

func (c *Client) CreateFormFieldPlacementCondition(d *FormFieldPlacementCondition) (*FormFieldPlacementCondition, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_placement_condition: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormFieldPlacementConditionRequestWithBody(c.Rootly.Server, d.FormFieldPlacementId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_field_placement_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPlacementCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_placement_condition: %s", err.Error())
	}

	return data.(*FormFieldPlacementCondition), nil
}

func (c *Client) GetFormFieldPlacementCondition(id string) (*FormFieldPlacementCondition, error) {
	req, err := rootlygo.NewGetFormFieldPlacementConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_field_placement_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPlacementCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_placement_condition: %s", err.Error())
	}

	return data.(*FormFieldPlacementCondition), nil
}

func (c *Client) UpdateFormFieldPlacementCondition(id string, form_field_placement_condition *FormFieldPlacementCondition) (*FormFieldPlacementCondition, error) {
	buffer, err := MarshalData(form_field_placement_condition)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_placement_condition: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormFieldPlacementConditionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_field_placement_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPlacementCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_placement_condition: %s", err.Error())
	}

	return data.(*FormFieldPlacementCondition), nil
}

func (c *Client) DeleteFormFieldPlacementCondition(id string) error {
	req, err := rootlygo.NewDeleteFormFieldPlacementConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_field_placement_condition: %s", err.Error())
	}

	return nil
}
