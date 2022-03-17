package client

import (
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type Service struct {
	ID                string `jsonapi:"primary,services"`
	Name              string `jsonapi:"attr,name,omitempty"`
	Slug              string `jsonapi:"attr,slug,omitempty"`
	Color             string `jsonapi:"attr,color,omitempty"`
	Description       string `jsonapi:"attr,description,omitempty"`
	PublicDescription string `jsonapi:"attr,public_description,omitempty"`
	//NotifyEmails  *[]string `json:"notify_emails,omitempty"`
	//SlackChannels *[]string `json:"slack_channels,omitempty"`
	//SlackAliases  *[]string `json:"slack_aliases,omitempty"`
}

func (c *Client) CreateService(s *Service) (*Service, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling service: %s", err.Error())
	}

	req, err := rootlygo.NewCreateServiceRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create service: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Service))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling service: %s", err.Error())
	}

	return data.(*Service), nil
}

func (c *Client) GetService(id string) (*Service, error) {
	req, err := rootlygo.NewGetServiceRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get service: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Service))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling service: %s", err.Error())
	}

	return data.(*Service), nil
}

func (c *Client) UpdateService(id string, s *Service) (*Service, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling service: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateServiceRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update service: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Service))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling service: %s", err.Error())
	}

	return data.(*Service), nil
}

func (c *Client) DeleteService(id string) error {
	req, err := rootlygo.NewDeleteServiceRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete service: %s", id)
	}

	return nil
}
