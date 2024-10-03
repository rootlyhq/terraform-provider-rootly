package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Workflow struct {
	ID string `jsonapi:"primary,workflows"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Command string `jsonapi:"attr,command,omitempty"`
  CommandFeedbackEnabled *bool `jsonapi:"attr,command_feedback_enabled,omitempty"`
  Wait string `jsonapi:"attr,wait,omitempty"`
  RepeatEveryDuration string `jsonapi:"attr,repeat_every_duration,omitempty"`
  RepeatOn []interface{} `jsonapi:"attr,repeat_on,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Locked *bool `jsonapi:"attr,locked,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  WorkflowGroupId string `jsonapi:"attr,workflow_group_id,omitempty"`
  TriggerParams map[string]interface{} `jsonapi:"attr,trigger_params,omitempty"`
  EnvironmentIds []interface{} `jsonapi:"attr,environment_ids,omitempty"`
  SeverityIds []interface{} `jsonapi:"attr,severity_ids,omitempty"`
  IncidentTypeIds []interface{} `jsonapi:"attr,incident_type_ids,omitempty"`
  IncidentRoleIds []interface{} `jsonapi:"attr,incident_role_ids,omitempty"`
  ServiceIds []interface{} `jsonapi:"attr,service_ids,omitempty"`
  FunctionalityIds []interface{} `jsonapi:"attr,functionality_ids,omitempty"`
  GroupIds []interface{} `jsonapi:"attr,group_ids,omitempty"`
  CauseIds []interface{} `jsonapi:"attr,cause_ids,omitempty"`
}

func (c *Client) ListWorkflows(params *rootlygo.ListWorkflowsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWorkflowsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	workflows, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Workflow)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return workflows, nil
}

func (c *Client) CreateWorkflow(d *Workflow) (*Workflow, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWorkflowRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create workflow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Workflow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow: %s", err.Error())
	}

	return data.(*Workflow), nil
}

func (c *Client) GetWorkflow(id string) (*Workflow, error) {
	req, err := rootlygo.NewGetWorkflowRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get workflow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Workflow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow: %s", err.Error())
	}

	return data.(*Workflow), nil
}

func (c *Client) UpdateWorkflow(id string, workflow *Workflow) (*Workflow, error) {
	buffer, err := MarshalData(workflow)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Workflow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow: %s", err.Error())
	}

	return data.(*Workflow), nil
}

func (c *Client) DeleteWorkflow(id string) error {
	req, err := rootlygo.NewDeleteWorkflowRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete workflow: %s", err.Error())
	}

	return nil
}
