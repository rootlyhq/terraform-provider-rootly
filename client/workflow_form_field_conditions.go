package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type WorkflowFormFieldCondition struct {
	ID string `jsonapi:"primary,workflow_form_field_conditions"`
	WorkflowId string `jsonapi:"attr,workflow_id,omitempty"`
  FormFieldId string `jsonapi:"attr,form_field_id,omitempty"`
  IncidentCondition string `jsonapi:"attr,incident_condition,omitempty"`
  Values []interface{} `jsonapi:"attr,values,omitempty"`
  SelectedCatalogEntityIds []interface{} `jsonapi:"attr,selected_catalog_entity_ids,omitempty"`
  SelectedFunctionalityIds []interface{} `jsonapi:"attr,selected_functionality_ids,omitempty"`
  SelectedGroupIds []interface{} `jsonapi:"attr,selected_group_ids,omitempty"`
  SelectedOptionIds []interface{} `jsonapi:"attr,selected_option_ids,omitempty"`
  SelectedServiceIds []interface{} `jsonapi:"attr,selected_service_ids,omitempty"`
  SelectedUserIds []interface{} `jsonapi:"attr,selected_user_ids,omitempty"`
}

func (c *Client) ListWorkflowFormFieldConditions(id string, params *rootlygo.ListWorkflowFormFieldConditionsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWorkflowFormFieldConditionsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	workflow_form_field_conditions, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(WorkflowFormFieldCondition)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return workflow_form_field_conditions, nil
}

func (c *Client) CreateWorkflowFormFieldCondition(d *WorkflowFormFieldCondition) (*WorkflowFormFieldCondition, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_form_field_condition: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWorkflowFormFieldConditionRequestWithBody(c.Rootly.Server, d.WorkflowId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create workflow_form_field_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowFormFieldCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_form_field_condition: %s", err.Error())
	}

	return data.(*WorkflowFormFieldCondition), nil
}

func (c *Client) GetWorkflowFormFieldCondition(id string) (*WorkflowFormFieldCondition, error) {
	req, err := rootlygo.NewGetWorkflowFormFieldConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get workflow_form_field_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowFormFieldCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_form_field_condition: %s", err.Error())
	}

	return data.(*WorkflowFormFieldCondition), nil
}

func (c *Client) UpdateWorkflowFormFieldCondition(id string, workflow_form_field_condition *WorkflowFormFieldCondition) (*WorkflowFormFieldCondition, error) {
	buffer, err := MarshalData(workflow_form_field_condition)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_form_field_condition: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowFormFieldConditionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow_form_field_condition: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowFormFieldCondition))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_form_field_condition: %s", err.Error())
	}

	return data.(*WorkflowFormFieldCondition), nil
}

func (c *Client) DeleteWorkflowFormFieldCondition(id string) error {
	req, err := rootlygo.NewDeleteWorkflowFormFieldConditionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete workflow_form_field_condition: %s", err.Error())
	}

	return nil
}
