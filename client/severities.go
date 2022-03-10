package client

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
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

func (s Severity) Marshal() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, s); err != nil {
		return nil, errors.Errorf("Error marshaling severity (creation): %s", err.Error())
	}

	return buffer, nil
}

func (c *Client) CreateSeverity(s *Severity) (*Severity, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, s); err != nil {
		return nil, errors.Errorf("Error marshaling severity (creation): %s", err.Error())
	}

	req, err := rootlygo.NewCreateSeverityRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (creation): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to create severity" + "\n\n" + buffer.String())
	}

	severity := new(Severity)
	if err := jsonapi.UnmarshalPayload(resp.Body, severity); err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (creation): %s", err.Error())
	}

	return severity, nil
}

func (c *Client) GetSeverity(id string) (*Severity, error) {
	req, err := rootlygo.NewGetSeverityRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get severity: %s", id)
	}

	severity := new(Severity)
	if err := jsonapi.UnmarshalPayload(resp.Body, severity); err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (read): %s", err.Error())
	}

	return severity, nil
}

func (c *Client) UpdateSeverity(id string, s *Severity) (*Severity, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, s); err != nil {
		return nil, errors.Errorf("Error marshaling severity (update): %s", err.Error())
	}

	req, err := rootlygo.NewUpdateSeverityRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (update): %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update severity: %s", id)
	}

	severity := new(Severity)
	if err := jsonapi.UnmarshalPayload(resp.Body, severity); err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (update): %s", err.Error())
	}

	return severity, nil
}

func (c *Client) DeleteSeverity(id string) error {
	req, err := rootlygo.NewDeleteSeverityRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error unmarshaling severity (delete): %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete severity: %s", id)
	}

	return nil
}
