package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/google/jsonapi"
)

const SeveritiesPath = "severities"

type Severity struct {
	ID          string `jsonapi:"primary,severities"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Slug        string `jsonapi:"attr,slug,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	Severity    string `jsonapi:"attr,severity,omitempty"`
	Color       string `jsonapi:"attr,color,omitempty"`
	//NotifyEmails  *[]string `json:"notify_emails,omitempty"`
	//SlackChannels *[]string `json:"slack_channels,omitempty"`
	//SlackAliases  *[]string `json:"slack_aliases,omitempty"`
}

func (c *Client) CreateSeverity(s *Severity) (*Severity, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, s); err != nil {
		return nil, errors.New("TODO")
	}
	req, err := http.NewRequest("POST", c.makeUrl(fmt.Sprintf("%s", SeveritiesPath)), buffer)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//severityResponse := SeverityResponse{}
	//err = json.Unmarshal(body, &severityResponse)
	//if err != nil {
	//	return nil, err
	//}

	severity := new(Severity)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), severity); err != nil {
		return nil, errors.New("TODO")
	}

	return severity, nil
}

func (c *Client) GetSeverity(id string) (*Severity, error) {
	var severity Severity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", SeveritiesPath, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &severity)
	if err != nil {
		return nil, err
	}

	return &severity, nil
}

func (c *Client) UpdateSeverity(id string, s *Severity) (*Severity, error) {
	data := &Data{
		Type: "severities",
		Attributes: map[string]interface{}{
			"name": "value",
		},
	}
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", SeveritiesPath, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	severity := Severity{}
	err = json.Unmarshal(body, &severity)
	if err != nil {
		return nil, err
	}

	return &severity, nil
}

func (c *Client) DeleteSeverity(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", SeveritiesPath, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
