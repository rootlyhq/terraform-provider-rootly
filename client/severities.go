package client

import (
	"reflect"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Severity struct {
	ID          string `jsonapi:"primary,severities"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Slug        string `jsonapi:"attr,slug,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Severity    string `jsonapi:"attr,severity,omitempty"`
	//NotifyEmails  *[]string `json:"notify_emails,omitempty"`
	//SlackChannels *[]string `json:"slack_channels,omitempty"`
	//SlackAliases  *[]string `json:"slack_aliases,omitempty"`
}

func (c *Client) ListSeverities(params *rootlygo.ListSeveritiesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListSeveritiesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	items, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Severity)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return items, nil
}

func (c *Client) CreateSeverity(s *Severity) (*Severity, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling severity: %s", err.Error())
	}

	req, err := rootlygo.NewCreateSeverityRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create severity: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Severity))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling severity: %s", err.Error())
	}

	return data.(*Severity), nil
}

func (c *Client) GetSeverity(id string) (*Severity, error) {
	req, err := rootlygo.NewGetSeverityRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get severity: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Severity))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling severity: %s", err.Error())
	}

	return data.(*Severity), nil
}

func (c *Client) UpdateSeverity(id string, s *Severity) (*Severity, error) {
	buffer, err := MarshalData(s)
	if err != nil {
		return nil, errors.Errorf("Error marshaling severity: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateSeverityRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update severity: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(Severity))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling severity: %s", err.Error())
	}

	return data.(*Severity), nil
}

func (c *Client) DeleteSeverity(id string) error {
	req, err := rootlygo.NewDeleteSeverityRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete severity: %s", id)
	}

	return nil
}
