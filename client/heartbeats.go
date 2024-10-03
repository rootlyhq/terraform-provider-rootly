package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Heartbeat struct {
	ID string `jsonapi:"primary,heartbeats"`
	Name string `jsonapi:"attr,name,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  AlertSummary string `jsonapi:"attr,alert_summary,omitempty"`
  AlertUrgencyId string `jsonapi:"attr,alert_urgency_id,omitempty"`
  Interval int `jsonapi:"attr,interval,omitempty"`
  IntervalUnit string `jsonapi:"attr,interval_unit,omitempty"`
  NotificationTargetId string `jsonapi:"attr,notification_target_id,omitempty"`
  NotificationTargetType string `jsonapi:"attr,notification_target_type,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Status string `jsonapi:"attr,status,omitempty"`
  LastPingedAt string `jsonapi:"attr,last_pinged_at,omitempty"`
  ExpiresAt string `jsonapi:"attr,expires_at,omitempty"`
}

func (c *Client) ListHeartbeats(params *rootlygo.ListHeartbeatsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListHeartbeatsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	heartbeats, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Heartbeat)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return heartbeats, nil
}

func (c *Client) CreateHeartbeat(d *Heartbeat) (*Heartbeat, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling heartbeat: %s", err.Error())
	}

	req, err := rootlygo.NewCreateHeartbeatRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create heartbeat: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Heartbeat))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling heartbeat: %s", err.Error())
	}

	return data.(*Heartbeat), nil
}

func (c *Client) GetHeartbeat(id string) (*Heartbeat, error) {
	req, err := rootlygo.NewGetHeartbeatRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get heartbeat: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Heartbeat))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling heartbeat: %s", err.Error())
	}

	return data.(*Heartbeat), nil
}

func (c *Client) UpdateHeartbeat(id string, heartbeat *Heartbeat) (*Heartbeat, error) {
	buffer, err := MarshalData(heartbeat)
	if err != nil {
		return nil, errors.Errorf("Error marshaling heartbeat: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateHeartbeatRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update heartbeat: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Heartbeat))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling heartbeat: %s", err.Error())
	}

	return data.(*Heartbeat), nil
}

func (c *Client) DeleteHeartbeat(id string) error {
	req, err := rootlygo.NewDeleteHeartbeatRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete heartbeat: %s", err.Error())
	}

	return nil
}
