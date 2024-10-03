package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RetrospectiveProcessGroup struct {
	ID string `jsonapi:"primary,retrospective_process_groups"`
	RetrospectiveProcessId string `jsonapi:"attr,retrospective_process_id,omitempty"`
  SubStatusId string `jsonapi:"attr,sub_status_id,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListRetrospectiveProcessGroups(id string, params *rootlygo.ListRetrospectiveProcessGroupsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveProcessGroupsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	retrospective_process_groups, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveProcessGroup)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return retrospective_process_groups, nil
}

func (c *Client) CreateRetrospectiveProcessGroup(d *RetrospectiveProcessGroup) (*RetrospectiveProcessGroup, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process_group: %s", err.Error())
	}

	req, err := rootlygo.NewCreateRetrospectiveProcessGroupRequestWithBody(c.Rootly.Server, d.RetrospectiveProcessId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create retrospective_process_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroup))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroup), nil
}

func (c *Client) GetRetrospectiveProcessGroup(id string) (*RetrospectiveProcessGroup, error) {
	req, err := rootlygo.NewGetRetrospectiveProcessGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get retrospective_process_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroup))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroup), nil
}

func (c *Client) UpdateRetrospectiveProcessGroup(id string, retrospective_process_group *RetrospectiveProcessGroup) (*RetrospectiveProcessGroup, error) {
	buffer, err := MarshalData(retrospective_process_group)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process_group: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRetrospectiveProcessGroupRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update retrospective_process_group: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcessGroup))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process_group: %s", err.Error())
	}

	return data.(*RetrospectiveProcessGroup), nil
}

func (c *Client) DeleteRetrospectiveProcessGroup(id string) error {
	req, err := rootlygo.NewDeleteRetrospectiveProcessGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete retrospective_process_group: %s", err.Error())
	}

	return nil
}
