package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type WorkflowTask struct {
	ID          string `jsonapi:"primary,workflow_tasks"`
	WorkflowId  string `jsonapi:"attr,workflow_id,omitempty"`
	TaskParams  map[string]interface{} `jsonapi:"attr,task_params,omitempty"`
}

func (c *Client) CreateWorkflowTask(i *WorkflowTask) (*WorkflowTask, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_task: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWorkflowTaskRequestWithBody(c.Rootly.Server, i.WorkflowId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create workflow_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_task: %s", err.Error())
	}

	return data.(*WorkflowTask), nil
}

func (c *Client) GetWorkflowTask(id string) (*WorkflowTask, error) {
	req, err := rootlygo.NewGetWorkflowTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get workflow_task: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_task: %s", err.Error())
	}

	return data.(*WorkflowTask), nil
}

func (c *Client) UpdateWorkflowTask(id string, i *WorkflowTask) (*WorkflowTask, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow_task: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowTaskRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow_task: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(WorkflowTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow_task: %s", err.Error())
	}

	return data.(*WorkflowTask), nil
}

func (c *Client) DeleteWorkflowTask(id string) error {
	req, err := rootlygo.NewDeleteWorkflowTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete workflow_task: %s", id)
	}

	return nil
}
