package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type SubStatus struct {
	ID string `jsonapi:"primary,sub_statuses"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  ParentStatus string `jsonapi:"attr,parent_status,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListSubStatuses(params *rootlygo.ListSubStatusesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListSubStatusesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	sub_statuses, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(SubStatus)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return sub_statuses, nil
}

func (c *Client) CreateSubStatus(d *SubStatus) (*SubStatus, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling sub_status: %s", err.Error())
	}

	req, err := rootlygo.NewCreateSubStatusRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(SubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling sub_status: %s", err.Error())
	}

	return data.(*SubStatus), nil
}

func (c *Client) GetSubStatus(id string) (*SubStatus, error) {
	req, err := rootlygo.NewGetSubStatusRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(SubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling sub_status: %s", err.Error())
	}

	return data.(*SubStatus), nil
}

func (c *Client) UpdateSubStatus(id string, sub_status *SubStatus) (*SubStatus, error) {
	buffer, err := MarshalData(sub_status)
	if err != nil {
		return nil, errors.Errorf("Error marshaling sub_status: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateSubStatusRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(SubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling sub_status: %s", err.Error())
	}

	return data.(*SubStatus), nil
}

func (c *Client) DeleteSubStatus(id string) error {
	req, err := rootlygo.NewDeleteSubStatusRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete sub_status: %s", err.Error())
	}

	return nil
}