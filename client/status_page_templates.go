package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type StatusPageTemplate struct {
	ID string `jsonapi:"primary,status_page_templates"`
	StatusPageId string `jsonapi:"attr,status_page_id,omitempty"`
  Title string `jsonapi:"attr,title,omitempty"`
  Body string `jsonapi:"attr,body,omitempty"`
  UpdateStatus string `jsonapi:"attr,update_status,omitempty"`
  ShouldNotifySubscribers *bool `jsonapi:"attr,should_notify_subscribers,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListStatusPageTemplates(id string, params *rootlygo.ListStatusPageTemplatesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListStatusPageTemplatesRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	status_page_templates, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(StatusPageTemplate)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return status_page_templates, nil
}

func (c *Client) CreateStatusPageTemplate(d *StatusPageTemplate) (*StatusPageTemplate, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling status_page_template: %s", err.Error())
	}

	req, err := rootlygo.NewCreateStatusPageTemplateRequestWithBody(c.Rootly.Server, d.StatusPageId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create status_page_template: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(StatusPageTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page_template: %s", err.Error())
	}

	return data.(*StatusPageTemplate), nil
}

func (c *Client) GetStatusPageTemplate(id string) (*StatusPageTemplate, error) {
	req, err := rootlygo.NewGetStatusPageTemplateRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get status_page_template: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(StatusPageTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page_template: %s", err.Error())
	}

	return data.(*StatusPageTemplate), nil
}

func (c *Client) UpdateStatusPageTemplate(id string, status_page_template *StatusPageTemplate) (*StatusPageTemplate, error) {
	buffer, err := MarshalData(status_page_template)
	if err != nil {
		return nil, errors.Errorf("Error marshaling status_page_template: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateStatusPageTemplateRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update status_page_template: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(StatusPageTemplate))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling status_page_template: %s", err.Error())
	}

	return data.(*StatusPageTemplate), nil
}

func (c *Client) DeleteStatusPageTemplate(id string) error {
	req, err := rootlygo.NewDeleteStatusPageTemplateRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete status_page_template: %s", id)
	}

	return nil
}
