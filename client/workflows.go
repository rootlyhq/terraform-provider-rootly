package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Workflow struct {
	ID          string `jsonapi:"primary,workflows"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Enabled     *bool  `jsonapi:"attr,enabled,omitempty"`
	WorkflowGroupId string `jsonapi:"attr,workflow_group_id,omitempty"`
	Position    int  `jsonapi:"attr,position,omitempty"`
	Command     string `jsonapi:"attr,command,omitempty"`
	TriggerParams map[string]interface{} `jsonapi:"attr,trigger_params,omitempty"`
	Wait        string `jsonapi:"attr,wait,omitempty"`
	RepeatEveryDuration string `jsonapi:"attr,repeat_every_duration,omitempty"`
	RepeatOn []interface{} `jsonapi:"attr,repeat_on,omitempty"`
	EnvironmentIds []interface{} `jsonapi:"attr,environment_ids,omitempty"`
	SeverityIds []interface{} `jsonapi:"attr,severity_ids,omitempty"`
	IncidentTypeIds []interface{} `jsonapi:"attr,incident_type_ids,omitempty"`
	ServiceIds []interface{} `jsonapi:"attr,service_ids,omitempty"`
	GroupIds []interface{} `jsonapi:"attr,group_ids,omitempty"`
}

func (c *Client) CreateWorkflow(i *Workflow) (*Workflow, error) {
	buffer, err := MarshalData(i)
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
		return nil, errors.Errorf("Failed to make request to get workflow: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Workflow))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling workflow: %s", err.Error())
	}

	return data.(*Workflow), nil
}

func (c *Client) UpdateWorkflow(id string, i *Workflow) (*Workflow, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling workflow: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWorkflowRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update workflow: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Workflow))
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
		return errors.Errorf("Failed to make request to delete workflow: %s", id)
	}

	return nil
}
