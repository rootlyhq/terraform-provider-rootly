package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type RetrospectiveStep struct {
	ID string `jsonapi:"primary,retrospective_steps"`
	Title string `jsonapi:"attr,title,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  IncidentRoleId string `jsonapi:"attr,incident_role_id,omitempty"`
  DueAfterDays int `jsonapi:"attr,due_after_days,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  Skippable *bool `jsonapi:"attr,skippable,omitempty"`
}

func (c *Client) ListRetrospectiveSteps(params *rootlygo.ListRetrospectiveStepsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveStepsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	retrospective_steps, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveStep)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return retrospective_steps, nil
}

func (c *Client) CreateRetrospectiveStep(d *RetrospectiveStep) (*RetrospectiveStep, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_step: %s", err.Error())
	}

	req, err := rootlygo.NewCreateRetrospectiveStepRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create retrospective_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_step: %s", err.Error())
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) GetRetrospectiveStep(id string) (*RetrospectiveStep, error) {
	req, err := rootlygo.NewGetRetrospectiveStepRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get retrospective_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_step: %s", err.Error())
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) UpdateRetrospectiveStep(id string, retrospective_step *RetrospectiveStep) (*RetrospectiveStep, error) {
	buffer, err := MarshalData(retrospective_step)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_step: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRetrospectiveStepRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update retrospective_step: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_step: %s", err.Error())
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) DeleteRetrospectiveStep(id string) error {
	req, err := rootlygo.NewDeleteRetrospectiveStepRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete retrospective_step: %s", err.Error())
	}

	return nil
}
