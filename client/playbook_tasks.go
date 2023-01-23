package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type PlaybookTask struct {
	ID string `jsonapi:"primary,playbook_tasks"`
	PlaybookId string `jsonapi:"attr,playbook_id,omitempty"`
  Task string `jsonapi:"attr,task,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
}

func (c *Client) ListPlaybookTasks(id string, params *rootlygo.ListPlaybookTasksParams) ([]interface{}, error) {
	req, err := rootlygo.NewListPlaybookTasksRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	playbook_tasks, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(PlaybookTask)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return playbook_tasks, nil
}

func (c *Client) CreatePlaybookTask(d *PlaybookTask) (*PlaybookTask, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling playbook_task: %s", err.Error())
	}

	req, err := rootlygo.NewCreatePlaybookTaskRequestWithBody(c.Rootly.Server, d.PlaybookId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create playbook_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PlaybookTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook_task: %s", err.Error())
	}

	return data.(*PlaybookTask), nil
}

func (c *Client) GetPlaybookTask(id string) (*PlaybookTask, error) {
	req, err := rootlygo.NewGetPlaybookTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get playbook_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PlaybookTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook_task: %s", err.Error())
	}

	return data.(*PlaybookTask), nil
}

func (c *Client) UpdatePlaybookTask(id string, playbook_task *PlaybookTask) (*PlaybookTask, error) {
	buffer, err := MarshalData(playbook_task)
	if err != nil {
		return nil, errors.Errorf("Error marshaling playbook_task: %s", err.Error())
	}

	req, err := rootlygo.NewUpdatePlaybookTaskRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update playbook_task: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(PlaybookTask))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook_task: %s", err.Error())
	}

	return data.(*PlaybookTask), nil
}

func (c *Client) DeletePlaybookTask(id string) error {
	req, err := rootlygo.NewDeletePlaybookTaskRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete playbook_task: %s", err.Error())
	}

	return nil
}
