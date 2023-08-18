package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Playbook struct {
	ID string `jsonapi:"primary,playbooks"`
	Title string `jsonapi:"attr,title,omitempty"`
  Summary string `jsonapi:"attr,summary,omitempty"`
  ExternalUrl string `jsonapi:"attr,external_url,omitempty"`
  SeverityIds []interface{} `jsonapi:"attr,severity_ids,omitempty"`
  EnvironmentIds []interface{} `jsonapi:"attr,environment_ids,omitempty"`
  FunctionalityIds []interface{} `jsonapi:"attr,functionality_ids,omitempty"`
  ServiceIds []interface{} `jsonapi:"attr,service_ids,omitempty"`
  GroupIds []interface{} `jsonapi:"attr,group_ids,omitempty"`
  IncidentTypeIds []interface{} `jsonapi:"attr,incident_type_ids,omitempty"`
}

func (c *Client) ListPlaybooks(params *rootlygo.ListPlaybooksParams) ([]interface{}, error) {
	req, err := rootlygo.NewListPlaybooksRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	playbooks, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Playbook)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return playbooks, nil
}

func (c *Client) CreatePlaybook(d *Playbook) (*Playbook, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling playbook: %s", err.Error())
	}

	req, err := rootlygo.NewCreatePlaybookRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create playbook: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Playbook))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook: %s", err.Error())
	}

	return data.(*Playbook), nil
}

func (c *Client) GetPlaybook(id string) (*Playbook, error) {
	req, err := rootlygo.NewGetPlaybookRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get playbook: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Playbook))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook: %s", err.Error())
	}

	return data.(*Playbook), nil
}

func (c *Client) UpdatePlaybook(id string, playbook *Playbook) (*Playbook, error) {
	buffer, err := MarshalData(playbook)
	if err != nil {
		return nil, errors.Errorf("Error marshaling playbook: %s", err.Error())
	}

	req, err := rootlygo.NewUpdatePlaybookRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update playbook: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Playbook))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling playbook: %s", err.Error())
	}

	return data.(*Playbook), nil
}

func (c *Client) DeletePlaybook(id string) error {
	req, err := rootlygo.NewDeletePlaybookRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete playbook: %s", err.Error())
	}

	return nil
}
