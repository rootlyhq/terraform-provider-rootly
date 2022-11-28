package client

import (
	
	
	"github.com/pkg/errors"
	
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

