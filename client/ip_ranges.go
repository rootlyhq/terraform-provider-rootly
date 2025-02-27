package client

import (
	"fmt"

	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type IpRanges struct {
	IntegrationsIpv4 []interface{} `jsonapi:"attr,integrations_ipv4,omitempty"`
	IntegrationsIpv6 []interface{} `jsonapi:"attr,integrations_ipv6,omitempty"`
	WebhooksIpv4     []interface{} `jsonapi:"attr,webhooks_ipv4,omitempty"`
	WebhooksIpv6     []interface{} `jsonapi:"attr,webhooks_ipv6,omitempty"`
}

func (c *Client) GetIpRanges() (*IpRanges, error) {
	req, err := rootlygo.NewGetIpRangesRequest(c.Rootly.Server)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get ip_ranges: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(IpRanges))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling ip_ranges: %w", err)
	}

	return data.(*IpRanges), nil
}
