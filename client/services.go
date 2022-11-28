package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Service struct {
	ID string `jsonapi:"primary,services"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  PublicDescription string `jsonapi:"attr,public_description,omitempty"`
  NotifyEmails []interface{} `jsonapi:"attr,notify_emails,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
  BackstageId string `jsonapi:"attr,backstage_id,omitempty"`
  PagerdutyId string `jsonapi:"attr,pagerduty_id,omitempty"`
  OpsgenieId string `jsonapi:"attr,opsgenie_id,omitempty"`
  GithubRepositoryName string `jsonapi:"attr,github_repository_name,omitempty"`
  GithubRepositoryBranch string `jsonapi:"attr,github_repository_branch,omitempty"`
  GitlabRepositoryName string `jsonapi:"attr,gitlab_repository_name,omitempty"`
  GitlabRepositoryBranch string `jsonapi:"attr,gitlab_repository_branch,omitempty"`
  EnvironmentIds []interface{} `jsonapi:"attr,environment_ids,omitempty"`
  ServiceIds []interface{} `jsonapi:"attr,service_ids,omitempty"`
  OwnersGroupIds []interface{} `jsonapi:"attr,owners_group_ids,omitempty"`
  SlackChannels []interface{} `jsonapi:"attr,slack_channels,omitempty"`
  SlackAliases []interface{} `jsonapi:"attr,slack_aliases,omitempty"`
}


func (c *Client) ListServices(params *rootlygo.ListServicesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListServicesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	services, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Service)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return services, nil
}


func (c *Client) CreateService(d *Service) (*Service, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling service: %s", err.Error())
	}

	req, err := rootlygo.NewCreateServiceRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create service: %s", err.Error())
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
		return nil, errors.Errorf("Failed to make request to get service: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Service))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling service: %s", err.Error())
	}

	return data.(*Service), nil
}


func (c *Client) UpdateService(id string, service *Service) (*Service, error) {
	buffer, err := MarshalData(service)
	if err != nil {
		return nil, errors.Errorf("Error marshaling service: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateServiceRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update service: %s", err.Error())
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
		return errors.Errorf("Failed to make request to delete service: %s", err.Error())
	}

	return nil
}

