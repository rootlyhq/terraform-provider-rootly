package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type RetrospectiveProcess struct {
	ID string `jsonapi:"primary,retrospective_processes"`
	CopyFrom string `jsonapi:"attr,copy_from,omitempty"`
	Name string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`
	IsDefault *bool `jsonapi:"attr,is_default,omitempty"`
	RetrospectiveProcessMatchingCriteria map[string]interface{} `jsonapi:"attr,retrospective_process_matching_criteria,omitempty"`
}

func (c *Client) ListRetrospectiveProcesses(params *rootlygo.ListRetrospectiveProcessesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRetrospectiveProcessesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	retrospective_processes, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(RetrospectiveProcess)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return retrospective_processes, nil
}

func (c *Client) CreateRetrospectiveProcess(d *RetrospectiveProcess) (*RetrospectiveProcess, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process: %s", err.Error())
	}

	req, err := rootlygo.NewCreateRetrospectiveProcessRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create retrospective_process: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcess))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process: %s", err.Error())
	}

	return data.(*RetrospectiveProcess), nil
}

func (c *Client) GetRetrospectiveProcess(id string) (*RetrospectiveProcess, error) {
	req, err := rootlygo.NewGetRetrospectiveProcessRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get retrospective_process: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcess))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process: %s", err.Error())
	}

	return data.(*RetrospectiveProcess), nil
}

func (c *Client) UpdateRetrospectiveProcess(id string, retrospective_process *RetrospectiveProcess) (*RetrospectiveProcess, error) {
	buffer, err := MarshalData(retrospective_process)
	if err != nil {
		return nil, errors.Errorf("Error marshaling retrospective_process: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRetrospectiveProcessRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update retrospective_process: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(RetrospectiveProcess))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling retrospective_process: %s", err.Error())
	}

	return data.(*RetrospectiveProcess), nil
}

func (c *Client) DeleteRetrospectiveProcess(id string) error {
	req, err := rootlygo.NewDeleteRetrospectiveProcessRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete retrospective_process: %s", err.Error())
	}

	return nil
}
