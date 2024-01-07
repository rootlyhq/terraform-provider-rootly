package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type WorkflowGroup struct {
	ID string `jsonapi:"primary,workflow_groups"`
	Kind string `jsonapi:"attr,kind,omitempty"`
  Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Icon string `jsonapi:"attr,icon,omitempty"`
  Expanded *bool `jsonapi:"attr,expanded,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListWorkflowGroups(params *rootlygo.ListWorkflowGroupsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWorkflowGroupsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	workflow_groups, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(WorkflowGroup)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return workflow_groups, nil
}

func (c *Client) CreateWorkflowGroup(d *WorkflowGroup) (*WorkflowGroup, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_group: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWorkflowGroupRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create workflow_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_group: %s", err.Error())
	}

	return data.(*WorkflowGroup), nil
}

func (c *Client) GetWorkflowGroup(id string) (*WorkflowGroup, error) {
	req, err := rootlygo.NewGetWorkflowGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get workflow_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_group: %s", err.Error())
	}

	return data.(*WorkflowGroup), nil
}

func (c *Client) UpdateWorkflowGroup(id string, workflow_group *WorkflowGroup) (*WorkflowGroup, error) {
	buffer, err := MarshalData(workflow_group)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_group: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowGroupRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_group: %s", err.Error())
	}

	return data.(*WorkflowGroup), nil
}

func (c *Client) DeleteWorkflowGroup(id string) error {
	req, err := rootlygo.NewDeleteWorkflowGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete workflow_group: %s", err.Error())
	}

	return nil
}
