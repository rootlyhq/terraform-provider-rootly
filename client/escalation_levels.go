package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type EscalationLevel struct {
	ID string `jsonapi:"primary,escalation_levels"`
	EscalationPolicyId string `jsonapi:"attr,escalation_policy_id,omitempty"`
  Delay int `jsonapi:"attr,delay,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  NotificationTargetParams []interface{} `jsonapi:"attr,notification_target_params,omitempty"`
}

func (c *Client) ListEscalationLevels(id string, params *rootlygo.ListEscalationLevelsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListEscalationLevelsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	escalation_levels, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(EscalationLevel)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return escalation_levels, nil
}

func (c *Client) CreateEscalationLevel(d *EscalationLevel) (*EscalationLevel, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling escalation_level: %s", err.Error())
	}

	req, err := rootlygo.NewCreateEscalationLevelRequestWithBody(c.Rootly.Server, d.EscalationPolicyId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create escalation_level: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationLevel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_level: %s", err.Error())
	}

	return data.(*EscalationLevel), nil
}

func (c *Client) GetEscalationLevel(id string) (*EscalationLevel, error) {
	req, err := rootlygo.NewGetEscalationLevelRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get escalation_level: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationLevel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_level: %s", err.Error())
	}

	return data.(*EscalationLevel), nil
}

func (c *Client) UpdateEscalationLevel(id string, escalation_level *EscalationLevel) (*EscalationLevel, error) {
	buffer, err := MarshalData(escalation_level)
	if err != nil {
		return nil, errors.Errorf("Error marshaling escalation_level: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateEscalationLevelRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update escalation_level: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationLevel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_level: %s", err.Error())
	}

	return data.(*EscalationLevel), nil
}

func (c *Client) DeleteEscalationLevel(id string) error {
	req, err := rootlygo.NewDeleteEscalationLevelRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete escalation_level: %s", err.Error())
	}

	return nil
}
