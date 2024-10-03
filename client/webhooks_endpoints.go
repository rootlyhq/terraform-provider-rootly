package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type WebhooksEndpoint struct {
	ID string `jsonapi:"primary,webhooks_endpoints"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Url string `jsonapi:"attr,url,omitempty"`
  EventTypes []interface{} `jsonapi:"attr,event_types,omitempty"`
  Secret string `jsonapi:"attr,secret,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
}

func (c *Client) ListWebhooksEndpoints(params *rootlygo.ListWebhooksEndpointsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWebhooksEndpointsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	webhooks_endpoints, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(WebhooksEndpoint)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return webhooks_endpoints, nil
}

func (c *Client) CreateWebhooksEndpoint(d *WebhooksEndpoint) (*WebhooksEndpoint, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling webhooks_endpoint: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWebhooksEndpointRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create webhooks_endpoint: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksEndpoint))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_endpoint: %s", err.Error())
	}

	return data.(*WebhooksEndpoint), nil
}

func (c *Client) GetWebhooksEndpoint(id string) (*WebhooksEndpoint, error) {
	req, err := rootlygo.NewGetWebhooksEndpointRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get webhooks_endpoint: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksEndpoint))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_endpoint: %s", err.Error())
	}

	return data.(*WebhooksEndpoint), nil
}

func (c *Client) UpdateWebhooksEndpoint(id string, webhooks_endpoint *WebhooksEndpoint) (*WebhooksEndpoint, error) {
	buffer, err := MarshalData(webhooks_endpoint)
	if err != nil {
		return nil, errors.Errorf("Error marshaling webhooks_endpoint: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWebhooksEndpointRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update webhooks_endpoint: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksEndpoint))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_endpoint: %s", err.Error())
	}

	return data.(*WebhooksEndpoint), nil
}

func (c *Client) DeleteWebhooksEndpoint(id string) error {
	req, err := rootlygo.NewDeleteWebhooksEndpointRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete webhooks_endpoint: %s", err.Error())
	}

	return nil
}
