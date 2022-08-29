package client

import (
	"reflect"
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type WorkflowCustomFieldSelection struct {
	ID string `jsonapi:"primary,workflow_custom_field_selections"`
	WorkflowId string `jsonapi:"attr,workflow_id,omitempty"`
	CustomFieldId int `jsonapi:"attr,custom_field_id,omitempty"`
	IncidentCondition string `jsonapi:"attr,incident_condition,omitempty"`
	Values []interface{} `jsonapi:"attr,values,omitempty"`
	SelectedOptionIds []interface{} `jsonapi:"attr,selected_option_ids,omitempty"`
}

func (c *Client) ListWorkflowCustomFieldSelections(workflowId string, params *rootlygo.ListWorkflowCustomFieldSelectionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWorkflowCustomFieldSelectionsRequest(c.Rootly.Server, workflowId, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	items, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(WorkflowCustomFieldSelection)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return items, nil
}

func (c *Client) CreateWorkflowCustomFieldSelection(workflowId string, i *WorkflowCustomFieldSelection) (*WorkflowCustomFieldSelection, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow custom field selection: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWorkflowCustomFieldSelectionRequestWithBody(c.Rootly.Server, workflowId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom field: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowCustomFieldSelection))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow custom field selection: %s", err.Error())
	}

	return data.(*WorkflowCustomFieldSelection), nil
}

func (c *Client) GetWorkflowCustomFieldSelection(id string) (*WorkflowCustomFieldSelection, error) {
	req, err := rootlygo.NewGetWorkflowCustomFieldSelectionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get custom field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowCustomFieldSelection))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field: %s", err.Error())
	}

	return data.(*WorkflowCustomFieldSelection), nil
}

func (c *Client) UpdateWorkflowCustomFieldSelection(id string, i *WorkflowCustomFieldSelection) (*WorkflowCustomFieldSelection, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom field: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowCustomFieldSelectionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom field: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowCustomFieldSelection))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field: %s", err.Error())
	}

	return data.(*WorkflowCustomFieldSelection), nil
}

func (c *Client) DeleteWorkflowCustomFieldSelection(id string) error {
	req, err := rootlygo.NewDeleteWorkflowCustomFieldSelectionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete custom field: %s", id)
	}

	return nil
}
