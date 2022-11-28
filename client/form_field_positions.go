package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type FormFieldPosition struct {
	ID string `jsonapi:"primary,form_field_positions"`
	FormFieldId string `jsonapi:"attr,form_field_id,omitempty"`
  Form string `jsonapi:"attr,form,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}


func (c *Client) ListFormFieldPositions(id string, params *rootlygo.ListFormFieldPositionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFormFieldPositionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	form_field_positions, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(FormFieldPosition)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return form_field_positions, nil
}


func (c *Client) CreateFormFieldPosition(d *FormFieldPosition) (*FormFieldPosition, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_position: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFormFieldPositionRequestWithBody(c.Rootly.Server, d.FormFieldId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create form_field_position: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPosition))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_position: %s", err.Error())
	}

	return data.(*FormFieldPosition), nil
}


func (c *Client) GetFormFieldPosition(id string) (*FormFieldPosition, error) {
	req, err := rootlygo.NewGetFormFieldPositionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get form_field_position: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPosition))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_position: %s", err.Error())
	}

	return data.(*FormFieldPosition), nil
}


func (c *Client) UpdateFormFieldPosition(id string, form_field_position *FormFieldPosition) (*FormFieldPosition, error) {
	buffer, err := MarshalData(form_field_position)
	if err != nil {
		return nil, errors.Errorf("Error marshaling form_field_position: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFormFieldPositionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update form_field_position: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(FormFieldPosition))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling form_field_position: %s", err.Error())
	}

	return data.(*FormFieldPosition), nil
}


func (c *Client) DeleteFormFieldPosition(id string) error {
	req, err := rootlygo.NewDeleteFormFieldPositionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete form_field_position: %s", err.Error())
	}

	return nil
}

