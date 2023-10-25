package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type IpRanges struct {
	Ipv4 []interface{} `jsonapi:"attr,ipv4,omitempty"`
	Ipv6 []interface{} `jsonapi:"attr,ipv4,omitempty"`
}

func (c *Client) GetIpRanges() (*IpRanges, error) {
	req, err := rootlygo.NewGetIpRangesRequest(c.Rootly.Server)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get ip_ranges: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IpRanges))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling ip_ranges: %s", err.Error())
	}

	return data.(*IpRanges), nil
}

