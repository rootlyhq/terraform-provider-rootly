package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type StatusPage struct {
	ID string `jsonapi:"primary,status_pages"`
	Title string `jsonapi:"attr,title,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  HeaderColor string `jsonapi:"attr,header_color,omitempty"`
  FooterColor string `jsonapi:"attr,footer_color,omitempty"`
  AllowSearchEngineIndex *bool `jsonapi:"attr,allow_search_engine_index,omitempty"`
  ShowUptime *bool `jsonapi:"attr,show_uptime,omitempty"`
  ShowUptimeLastDays int `jsonapi:"attr,show_uptime_last_days,omitempty"`
  Public *bool `jsonapi:"attr,public,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
}


func (c *Client) ListStatusPages(params *rootlygo.ListStatusPagesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListStatusPagesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	status_pages, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(StatusPage)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return status_pages, nil
}


func (c *Client) CreateStatusPage(d *StatusPage) (*StatusPage, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling status_page: %s", err.Error())
	}

	req, err := rootlygo.NewCreateStatusPageRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create status_page: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(StatusPage))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page: %s", err.Error())
	}

	return data.(*StatusPage), nil
}


func (c *Client) GetStatusPage(id string) (*StatusPage, error) {
	req, err := rootlygo.NewGetStatusPageRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get status_page: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(StatusPage))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page: %s", err.Error())
	}

	return data.(*StatusPage), nil
}


func (c *Client) UpdateStatusPage(id string, status_page *StatusPage) (*StatusPage, error) {
	buffer, err := MarshalData(status_page)
	if err != nil {
		return nil, errors.Errorf("Error marshaling status_page: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateStatusPageRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update status_page: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(StatusPage))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page: %s", err.Error())
	}

	return data.(*StatusPage), nil
}


func (c *Client) DeleteStatusPage(id string) error {
	req, err := rootlygo.NewDeleteStatusPageRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete status_page: %s", err.Error())
	}

	return nil
}

