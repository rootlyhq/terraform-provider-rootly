package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type AlertUrgency struct {
	ID string `jsonapi:"primary,alert_urgencies"`
	Name string `jsonapi:"attr,name,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListAlertUrgencies(params *rootlygo.ListAlertUrgenciesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListAlertUrgenciesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	alert_urgencies, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(AlertUrgency)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return alert_urgencies, nil
}

func (c *Client) CreateAlertUrgency(d *AlertUrgency) (*AlertUrgency, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling alert_urgency: %s", err.Error())
	}

	req, err := rootlygo.NewCreateAlertUrgencyRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create alert_urgency: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertUrgency))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_urgency: %s", err.Error())
	}

	return data.(*AlertUrgency), nil
}

func (c *Client) GetAlertUrgency(id string) (*AlertUrgency, error) {
	req, err := rootlygo.NewGetAlertUrgencyRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get alert_urgency: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertUrgency))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_urgency: %s", err.Error())
	}

	return data.(*AlertUrgency), nil
}

func (c *Client) UpdateAlertUrgency(id string, alert_urgency *AlertUrgency) (*AlertUrgency, error) {
	buffer, err := MarshalData(alert_urgency)
	if err != nil {
		return nil, errors.Errorf("Error marshaling alert_urgency: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateAlertUrgencyRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update alert_urgency: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertUrgency))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_urgency: %s", err.Error())
	}

	return data.(*AlertUrgency), nil
}

func (c *Client) DeleteAlertUrgency(id string) error {
	req, err := rootlygo.NewDeleteAlertUrgencyRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete alert_urgency: %s", err.Error())
	}

	return nil
}