package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type WorkflowGroup struct {
	ID          string `jsonapi:"primary,workflow_groups"`
	Kind        string `jsonapi:"attr,kind,omitempty"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Expanded    *bool  `jsonapi:"attr,expanded,omitempty"`
	Position    int  `jsonapi:"attr,position,omitempty"`
}

func (c *Client) CreateWorkflowGroup(i *WorkflowGroup) (*WorkflowGroup, error) {
	buffer, err := MarshalData(i)
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
		return nil, errors.Errorf("Failed to make request to get workflow_group: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_group: %s", err.Error())
	}

	return data.(*WorkflowGroup), nil
}

func (c *Client) UpdateWorkflowGroup(id string, i *WorkflowGroup) (*WorkflowGroup, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_group: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowGroupRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow_group: %s", id)
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
		return errors.Errorf("Failed to make request to delete workflow_group: %s", id)
	}

	return nil
}
