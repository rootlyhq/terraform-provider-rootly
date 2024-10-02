package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RetrospectiveProcessGroupStep struct {
	ID string `jsonapi:"primary,retrospective_process_group_steps"`
	RetrospectiveProcessGroupId string `jsonapi:"attr,retrospective_process_group_id,omitempty"`
  RetrospectiveStepId string `jsonapi:"attr,retrospective_step_id,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListRetrospectiveProcessGroupSteps(id string, params *rootlygo.ListRetrospectiveProcessGroupStepsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveProcessGroupStepsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	retrospective_process_group_steps, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveProcessGroupStep)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return retrospective_process_group_steps, nil
}

func (c *Client) CreateRetrospectiveProcessGroupStep(d *RetrospectiveProcessGroupStep) (*RetrospectiveProcessGroupStep, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process_group_step: %s", err.Error())
	}

	req, err := rootlygo.NewCreateRetrospectiveProcessGroupStepRequestWithBody(c.Rootly.Server, d.RetrospectiveProcessGroupId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create retrospective_process_group_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroupStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group_step: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroupStep), nil
}

func (c *Client) GetRetrospectiveProcessGroupStep(id string) (*RetrospectiveProcessGroupStep, error) {
	req, err := rootlygo.NewGetRetrospectiveProcessGroupStepRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get retrospective_process_group_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroupStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group_step: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroupStep), nil
}

func (c *Client) UpdateRetrospectiveProcessGroupStep(id string, retrospective_process_group_step *RetrospectiveProcessGroupStep) (*RetrospectiveProcessGroupStep, error) {
	buffer, err := MarshalData(retrospective_process_group_step)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process_group_step: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRetrospectiveProcessGroupStepRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update retrospective_process_group_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroupStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group_step: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroupStep), nil
}

func (c *Client) DeleteRetrospectiveProcessGroupStep(id string) error {
	req, err := rootlygo.NewDeleteRetrospectiveProcessGroupStepRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete retrospective_process_group_step: %s", err.Error())
	}

	return nil
}
