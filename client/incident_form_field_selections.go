package client

import (
	"reflect"

	"fmt"

	"github.com/google/jsonapi"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type IncidentFormFieldSelection struct {
	ID                string        `jsonapi:"primary,incident_form_field_selections"`
	IncidentId        string        `jsonapi:"attr,incident_id,omitempty"`
	FormFieldId       string        `jsonapi:"attr,form_field_id,omitempty"`
	Value             string        `jsonapi:"attr,value,omitempty"`
	SelectedOptionIds []interface{} `jsonapi:"attr,selected_option_ids,omitempty"`
	SelectedUserIds   []interface{} `jsonapi:"attr,selected_user_ids,omitempty"`
}

func (c *Client) ListIncidentFormFieldSelections(id rootlygo_.ID, params *rootlygo.ListIncidentFormFieldSelectionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentFormFieldSelectionsRequest(c.Rootly.Server, id.String(), params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	incident_form_field_selections, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentFormFieldSelection)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return incident_form_field_selections, nil
}

func (c *Client) CreateIncidentFormFieldSelection(d *IncidentFormFieldSelection) (*IncidentFormFieldSelection, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling incident_form_field_selection: %w", err)
	}

	req, err := rootlygo.NewCreateIncidentFormFieldSelectionRequestWithBody(c.Rootly.Server, d.IncidentId, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create incident_form_field_selection: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentFormFieldSelection))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling incident_form_field_selection: %w", err)
	}

	return data.(*IncidentFormFieldSelection), nil
}

func (c *Client) GetIncidentFormFieldSelection(id rootlygo_.ID) (*IncidentFormFieldSelection, error) {
	req, err := rootlygo.NewGetIncidentFormFieldSelectionRequest(c.Rootly.Server, id.String())
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get incident_form_field_selection: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentFormFieldSelection))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling incident_form_field_selection: %w", err)
	}

	return data.(*IncidentFormFieldSelection), nil
}

func (c *Client) UpdateIncidentFormFieldSelection(id rootlygo_.ID, incident_form_field_selection *IncidentFormFieldSelection) (*IncidentFormFieldSelection, error) {
	buffer, err := MarshalData(incident_form_field_selection)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling incident_form_field_selection: %w", err)
	}

	req, err := rootlygo.NewUpdateIncidentFormFieldSelectionRequestWithBody(c.Rootly.Server, id.String(), c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update incident_form_field_selection: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(IncidentFormFieldSelection))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling incident_form_field_selection: %w", err)
	}

	return data.(*IncidentFormFieldSelection), nil
}

func (c *Client) DeleteIncidentFormFieldSelection(id rootlygo_.ID) error {
	req, err := rootlygo.NewDeleteIncidentFormFieldSelectionRequest(c.Rootly.Server, id.String())
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete incident_form_field_selection: %w", err)
	}

	return nil
}
