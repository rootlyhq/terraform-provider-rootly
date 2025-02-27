package client

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Dashboard struct {
	ID     string `jsonapi:"primary,dashboards"`
	Slug   string `jsonapi:"attr,slug,omitempty"`
	Name   string `jsonapi:"attr,name,omitempty"`
	Owner  string `jsonapi:"attr,owner,omitempty"`
	UserId int    `jsonapi:"attr,user_id,omitempty"`
	Public *bool  `jsonapi:"attr,public,omitempty"`
}

func (c *Client) ListDashboards(params *rootlygo.ListDashboardsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListDashboardsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	dashboards, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Dashboard)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return dashboards, nil
}

func (c *Client) CreateDashboard(dashboard *Dashboard) (*Dashboard, error) {
	buffer, err := MarshalData(dashboard)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling dashboard: %w", err)
	}

	req, err := rootlygo.NewCreateDashboardRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create dashboard: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling dashboard during create: %w", err)
	}

	return data.(*Dashboard), nil
}

func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	req, err := rootlygo.NewGetDashboardRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get dashboard: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling dashboard: %w", err)
	}

	return data.(*Dashboard), nil
}

func (c *Client) UpdateDashboard(id string, dashboard *Dashboard) (*Dashboard, error) {
	buffer, err := MarshalData(dashboard)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling dashboard: %w", err)
	}

	req, err := rootlygo.NewUpdateDashboardRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update dashboard: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Dashboard))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling dashboard: %w", err)
	}

	return data.(*Dashboard), nil
}

func (c *Client) DeleteDashboard(id string) error {
	req, err := rootlygo.NewDeleteDashboardRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete dashboard: %w", err)
	}

	return nil
}
