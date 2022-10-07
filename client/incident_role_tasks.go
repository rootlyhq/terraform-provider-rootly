package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type IncidentRoleTask struct {
	ID string `jsonapi:"primary,incident_role_tasks"`
	IncidentRoleId string `jsonapi:"attr,incident_role_id,omitempty"`
  Task string `jsonapi:"attr,task,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Priority string `jsonapi:"attr,priority,omitempty"`
}

func (c *Client) ListIncidentRoleTasks(id string, params *rootlygo.ListIncidentRoleTasksParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentRoleTasksRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_role_tasks, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentRoleTask)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_role_tasks, nil
}

func (c *Client) CreateIncidentRoleTask(d *IncidentRoleTask) (*IncidentRoleTask, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_role_task: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentRoleTaskRequestWithBody(c.Rootly.Server, d.IncidentRoleId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_role_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRoleTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role_task: %s", err.Error())
	}

	return data.(*IncidentRoleTask), nil
}

func (c *Client) GetIncidentRoleTask(id string) (*IncidentRoleTask, error) {
	req, err := rootlygo.NewGetIncidentRoleTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_role_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRoleTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role_task: %s", err.Error())
	}

	return data.(*IncidentRoleTask), nil
}

func (c *Client) UpdateIncidentRoleTask(id string, incident_role_task *IncidentRoleTask) (*IncidentRoleTask, error) {
	buffer, err := MarshalData(incident_role_task)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_role_task: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentRoleTaskRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_role_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentRoleTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_role_task: %s", err.Error())
	}

	return data.(*IncidentRoleTask), nil
}

func (c *Client) DeleteIncidentRoleTask(id string) error {
	req, err := rootlygo.NewDeleteIncidentRoleTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_role_task: %s", err.Error())
	}

	return nil
}
