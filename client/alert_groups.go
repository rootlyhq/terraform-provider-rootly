package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type AlertGroup struct {
	ID string `jsonapi:"primary,alert_groups"`
	Name string `jsonapi:"attr,name,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  ConditionType string `jsonapi:"attr,condition_type,omitempty"`
  TimeWindow int `jsonapi:"attr,time_window,omitempty"`
  GroupByAlertTitle *bool `jsonapi:"attr,group_by_alert_title,omitempty"`
  GroupByAlertUrgency *bool `jsonapi:"attr,group_by_alert_urgency,omitempty"`
  DeletedAt string `jsonapi:"attr,deleted_at,omitempty"`
}

func (c *Client) ListAlertGroups(params *rootlygo.ListAlertGroupsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListAlertGroupsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	alert_groups, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(AlertGroup)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return alert_groups, nil
}

func (c *Client) CreateAlertGroup(d *AlertGroup) (*AlertGroup, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling alert_group: %s", err.Error())
	}

	req, err := rootlygo.NewCreateAlertGroupRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create alert_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_group: %s", err.Error())
	}

	return data.(*AlertGroup), nil
}

func (c *Client) GetAlertGroup(id string) (*AlertGroup, error) {
	req, err := rootlygo.NewGetAlertGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get alert_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_group: %s", err.Error())
	}

	return data.(*AlertGroup), nil
}

func (c *Client) UpdateAlertGroup(id string, alert_group *AlertGroup) (*AlertGroup, error) {
	buffer, err := MarshalData(alert_group)
	if err != nil {
		return nil, errors.Errorf("Error marshaling alert_group: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateAlertGroupRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update alert_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(AlertGroup))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling alert_group: %s", err.Error())
	}

	return data.(*AlertGroup), nil
}

func (c *Client) DeleteAlertGroup(id string) error {
	req, err := rootlygo.NewDeleteAlertGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete alert_group: %s", err.Error())
	}

	return nil
}