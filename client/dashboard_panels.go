package client

import (
	"reflect"
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type DashboardPanel struct {
	ID          string `jsonapi:"primary,dashboard_panels"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Params      map[string]interface{} `jsonapi:"attr,params,omitempty"`
	Position    map[string]interface{} `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListDashboardPanels(dashboardId string, params *rootlygo.ListDashboardPanelsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListDashboardPanelsRequest(c.Rootly.Server, dashboardId, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	dashboard_panels, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(DashboardPanel)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return dashboard_panels, nil
}

func (c *Client) CreateDashboardPanel(dashboardId string, dashboard_panel *DashboardPanel) (*DashboardPanel, error) {
	buffer, err := MarshalData(dashboard_panel)
	if err != nil {
		return nil, errors.Errorf("Error marshaling dashboard_panel: %s", err.Error())
	}

	req, err := rootlygo.NewCreateDashboardPanelRequestWithBody(c.Rootly.Server, dashboardId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create dashboard_panel: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(DashboardPanel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard_panel: %s", err.Error())
	}

	return data.(*DashboardPanel), nil
}

func (c *Client) GetDashboardPanel(id string, params *rootlygo.GetDashboardPanelParams) (*DashboardPanel, error) {
	req, err := rootlygo.NewGetDashboardPanelRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get dashboard_panel: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(DashboardPanel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard_panel: %s", err.Error())
	}

	return data.(*DashboardPanel), nil
}

func (c *Client) UpdateDashboardPanel(id string, dashboard_panel *DashboardPanel) (*DashboardPanel, error) {
	buffer, err := MarshalData(dashboard_panel)
	if err != nil {
		return nil, errors.Errorf("Error marshaling dashboard_panel: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateDashboardPanelRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update dashboard_panel: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(DashboardPanel))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling dashboard_panel: %s", err.Error())
	}

	return data.(*DashboardPanel), nil
}

func (c *Client) DeleteDashboardPanel(id string) error {
	req, err := rootlygo.NewDeleteDashboardPanelRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete dashboard_panel: %s", id)
	}

	return nil
}
