package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type WebhooksDelivery struct {
	ID string `jsonapi:"primary,webhooks_deliveries"`
	EndpointId string `jsonapi:"attr,endpoint_id,omitempty"`
  Delivered *bool `jsonapi:"attr,delivered,omitempty"`
  Attempts int `jsonapi:"attr,attempts,omitempty"`
  Payload string `jsonapi:"attr,payload,omitempty"`
  LastAttemptAt string `jsonapi:"attr,last_attempt_at,omitempty"`
}

func (c *Client) ListWebhooksDeliveries(params *rootlygo.ListWebhooksDeliveriesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListWebhooksDeliveriesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	webhooks_deliveries, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(WebhooksDelivery)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return webhooks_deliveries, nil
}

func (c *Client) CreateWebhooksDelivery(d *WebhooksDelivery) (*WebhooksDelivery, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling webhooks_delivery: %s", err.Error())
	}

	req, err := rootlygo.NewCreateWebhooksDeliveryRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create webhooks_delivery: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksDelivery))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_delivery: %s", err.Error())
	}

	return data.(*WebhooksDelivery), nil
}

func (c *Client) GetWebhooksDelivery(id string) (*WebhooksDelivery, error) {
	req, err := rootlygo.NewGetWebhooksDeliveryRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get webhooks_delivery: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksDelivery))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_delivery: %s", err.Error())
	}

	return data.(*WebhooksDelivery), nil
}

func (c *Client) UpdateWebhooksDelivery(id string, webhooks_delivery *WebhooksDelivery) (*WebhooksDelivery, error) {
	buffer, err := MarshalData(webhooks_delivery)
	if err != nil {
		return nil, errors.Errorf("Error marshaling webhooks_delivery: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateWebhooksDeliveryRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update webhooks_delivery: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(WebhooksDelivery))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling webhooks_delivery: %s", err.Error())
	}

	return data.(*WebhooksDelivery), nil
}

func (c *Client) DeleteWebhooksDelivery(id string) error {
	req, err := rootlygo.NewDeleteWebhooksDeliveryRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete webhooks_delivery: %s", err.Error())
	}

	return nil
}
