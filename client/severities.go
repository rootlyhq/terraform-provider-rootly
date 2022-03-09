package client

import (
	"bytes"
	"fmt"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

const SeveritiesPath = "/v1/severities"

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
		fmt.Println(err)
		return nil, errors.Errorf("Failed to make request to create severity" + "\n\n" + buffer.String())
	}

	severity := new(Severity)
	if err := jsonapi.UnmarshalPayload(resp.Body, severity); err != nil {
		return nil, errors.Errorf("Error unmarshaling severity (creation): %s", err.Error())
	}

	return severity, nil
}

//func (c *Client) GetSeverity(id string) (*Severity, error) {
//	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", SeveritiesPath, id), nil)
//	if err != nil {
//		return nil, err
//	}
//
//	body, err := c.doRequest(req)
//	if err != nil {
//		return nil, err
//	}
//
//	severity := new(Severity)
//	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), severity); err != nil {
//		return nil, errors.Errorf("Error unmarshaling severity (creation): %s", err.Error())
//	}
//
//	return severity, nil
//}
//
//func (c *Client) UpdateSeverity(id string, s *Severity) (*Severity, error) {
//	buffer := new(bytes.Buffer)
//	if err := jsonapi.MarshalPayload(buffer, s); err != nil {
//		return nil, errors.Errorf("Error marshaling severity (update): %s", err.Error())
//	}
//	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", SeveritiesPath, id), buffer)
//	if err != nil {
//		return nil, err
//	}
//
//	body, err := c.doRequest(req)
//	if err != nil {
//		return nil, err
//	}
//
//	severity := new(Severity)
//	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), severity); err != nil {
//		return nil, errors.Errorf("Error unmarshaling severity (creation): %s", err.Error())
//	}
//
//	return severity, nil
//}
//
//func (c *Client) DeleteSeverity(id string) error {
//	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", SeveritiesPath, id), nil)
//	if err != nil {
//		return err
//	}
//
//	_, err = c.doRequest(req)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
