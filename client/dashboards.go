package client

import (
	"reflect"
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Dashboard struct {
	ID          string `jsonapi:"primary,dashboards"`
	Slug        string `jsonapi:"attr,slug,omitempty"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Owner       string `jsonapi:"attr,owner,omitempty"`
	UserId      int `jsonapi:"attr,user_id,omitempty"`
	Public      *bool `jsonapi:"attr,public,omitempty"`
}

func (c *Client) ListDashboards(params *rootlygo.ListDashboardsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListDashboardsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	dashboards, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Dashboard)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return dashboards, nil
}

func (c *Client) CreateDashboard(dashboard *Dashboard) (*Dashboard, error) {
	buffer, err := MarshalData(dashboard)
	if err != nil {
		return nil, errors.Errorf("Error marshaling dashboard: %s", err.Error())
	}

	req, err := rootlygo.NewCreateDashboardRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create dashboard: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard during create: %s", err.Error())
	}

	return data.(*Dashboard), nil
}

func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	req, err := rootlygo.NewGetDashboardRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get dashboard: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard: %s", err.Error())
	}

	return data.(*Dashboard), nil
}

func (c *Client) UpdateDashboard(id string, dashboard *Dashboard) (*Dashboard, error) {
	buffer, err := MarshalData(dashboard)
	if err != nil {
		return nil, errors.Errorf("Error marshaling dashboard: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateDashboardRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update dashboard: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard: %s", err.Error())
	}

	return data.(*Dashboard), nil
}

func (c *Client) DeleteDashboard(id string) error {
	req, err := rootlygo.NewDeleteDashboardRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete dashboard: %s", id)
	}

	return nil
}
