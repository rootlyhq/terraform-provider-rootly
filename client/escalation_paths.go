package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type EscalationPath struct {
	ID string `jsonapi:"primary,escalation_paths"`
	Name string `jsonapi:"attr,name,omitempty"`
  Default *bool `jsonapi:"attr,default,omitempty"`
  NotificationType string `jsonapi:"attr,notification_type,omitempty"`
  EscalationPolicyId string `jsonapi:"attr,escalation_policy_id,omitempty"`
  Repeat *bool `jsonapi:"attr,repeat,omitempty"`
  RepeatCount int `jsonapi:"attr,repeat_count,omitempty"`
  Rules []interface{} `jsonapi:"attr,rules,omitempty"`
}

func (c *Client) ListEscalationPaths(id string, params *rootlygo.ListEscalationPathsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListEscalationPathsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	escalation_paths, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(EscalationPath)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return escalation_paths, nil
}

func (c *Client) CreateEscalationPath(d *EscalationPath) (*EscalationPath, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling escalation_path: %s", err.Error())
	}

	req, err := rootlygo.NewCreateEscalationPathRequestWithBody(c.Rootly.Server, d.EscalationPolicyId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create escalation_path: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationPath))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_path: %s", err.Error())
	}

	return data.(*EscalationPath), nil
}

func (c *Client) GetEscalationPath(id string) (*EscalationPath, error) {
	req, err := rootlygo.NewGetEscalationPathRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get escalation_path: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationPath))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_path: %s", err.Error())
	}

	return data.(*EscalationPath), nil
}

func (c *Client) UpdateEscalationPath(id string, escalation_path *EscalationPath) (*EscalationPath, error) {
	buffer, err := MarshalData(escalation_path)
	if err != nil {
		return nil, errors.Errorf("Error marshaling escalation_path: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateEscalationPathRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update escalation_path: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(EscalationPath))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling escalation_path: %s", err.Error())
	}

	return data.(*EscalationPath), nil
}

func (c *Client) DeleteEscalationPath(id string) error {
	req, err := rootlygo.NewDeleteEscalationPathRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete escalation_path: %s", err.Error())
	}

	return nil
}
