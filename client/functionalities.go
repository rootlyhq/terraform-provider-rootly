package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Functionality struct {
	ID string `jsonapi:"primary,functionalities"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  PublicDescription string `jsonapi:"attr,public_description,omitempty"`
  NotifyEmails []interface{} `jsonapi:"attr,notify_emails,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  EnvironmentIds []interface{} `jsonapi:"attr,environment_ids,omitempty"`
  ServiceIds []interface{} `jsonapi:"attr,service_ids,omitempty"`
  OwnersGroupIds []interface{} `jsonapi:"attr,owners_group_ids,omitempty"`
  OwnersUserIds []interface{} `jsonapi:"attr,owners_user_ids,omitempty"`
  SlackChannels []interface{} `jsonapi:"attr,slack_channels,omitempty"`
  SlackAliases []interface{} `jsonapi:"attr,slack_aliases,omitempty"`
}

func (c *Client) ListFunctionalities(params *rootlygo.ListFunctionalitiesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListFunctionalitiesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	functionalities, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Functionality)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return functionalities, nil
}

func (c *Client) CreateFunctionality(d *Functionality) (*Functionality, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling functionality: %s", err.Error())
	}

	req, err := rootlygo.NewCreateFunctionalityRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create functionality: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) GetFunctionality(id string) (*Functionality, error) {
	req, err := rootlygo.NewGetFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get functionality: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) UpdateFunctionality(id string, functionality *Functionality) (*Functionality, error) {
	buffer, err := MarshalData(functionality)
	if err != nil {
		return nil, errors.Errorf("Error marshaling functionality: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateFunctionalityRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update functionality: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Functionality))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling functionality: %s", err.Error())
	}

	return data.(*Functionality), nil
}

func (c *Client) DeleteFunctionality(id string) error {
	req, err := rootlygo.NewDeleteFunctionalityRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete functionality: %s", err.Error())
	}

	return nil
}
