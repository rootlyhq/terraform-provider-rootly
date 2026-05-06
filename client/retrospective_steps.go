package client

import (
	"reflect"

	"fmt"

	"github.com/google/jsonapi"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type RetrospectiveStep struct {
	ID                     string `jsonapi:"primary,retrospective_steps"`
	RetrospectiveProcessId string `jsonapi:"attr,retrospective_process_id,omitempty"`
	Title                  string `jsonapi:"attr,title,omitempty"`
	Slug                   string `jsonapi:"attr,slug,omitempty"`
	Description            string `jsonapi:"attr,description,omitempty"`
	IncidentRoleId         string `jsonapi:"attr,incident_role_id,omitempty"`
	DueAfterDays           int    `jsonapi:"attr,due_after_days,omitempty"`
	Position               int    `jsonapi:"attr,position,omitempty"`
	Skippable              *bool  `jsonapi:"attr,skippable,omitempty"`
}

func (c *Client) ListRetrospectiveSteps(id rootlygo_.ID, params *rootlygo.ListRetrospectiveStepsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveStepsRequest(c.Rootly.Server, id.String(), params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	retrospective_steps, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveStep)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return retrospective_steps, nil
}

func (c *Client) CreateRetrospectiveStep(d *RetrospectiveStep) (*RetrospectiveStep, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling retrospective_step: %w", err)
	}

	req, err := rootlygo.NewCreateRetrospectiveStepRequestWithBody(c.Rootly.Server, d.RetrospectiveProcessId, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create retrospective_step: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling retrospective_step: %w", err)
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) GetRetrospectiveStep(id rootlygo_.ID) (*RetrospectiveStep, error) {
	req, err := rootlygo.NewGetRetrospectiveStepRequest(c.Rootly.Server, id.String())
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get retrospective_step: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling retrospective_step: %w", err)
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) UpdateRetrospectiveStep(id rootlygo_.ID, retrospective_step *RetrospectiveStep) (*RetrospectiveStep, error) {
	buffer, err := MarshalData(retrospective_step)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling retrospective_step: %w", err)
	}

	req, err := rootlygo.NewUpdateRetrospectiveStepRequestWithBody(c.Rootly.Server, id.String(), c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update retrospective_step: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveStep))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling retrospective_step: %w", err)
	}

	return data.(*RetrospectiveStep), nil
}

func (c *Client) DeleteRetrospectiveStep(id rootlygo_.ID) error {
	req, err := rootlygo.NewDeleteRetrospectiveStepRequest(c.Rootly.Server, id.String())
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete retrospective_step: %w", err)
	}

	return nil
}
