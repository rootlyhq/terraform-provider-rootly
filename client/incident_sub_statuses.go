package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type IncidentSubStatus struct {
	ID string `jsonapi:"primary,incident_sub_statuses"`
	IncidentId string `jsonapi:"attr,incident_id,omitempty"`
  SubStatusId string `jsonapi:"attr,sub_status_id,omitempty"`
  AssignedAt string `jsonapi:"attr,assigned_at,omitempty"`
  AssignedByUserId int `jsonapi:"attr,assigned_by_user_id,omitempty"`
}

func (c *Client) ListIncidentSubStatuses(id string, params *rootlygo.ListIncidentSubStatusesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentSubStatusesRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_sub_statuses, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentSubStatus)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_sub_statuses, nil
}

func (c *Client) CreateIncidentSubStatus(d *IncidentSubStatus) (*IncidentSubStatus, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_sub_status: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentSubStatusRequestWithBody(c.Rootly.Server, d.IncidentId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentSubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_sub_status: %s", err.Error())
	}

	return data.(*IncidentSubStatus), nil
}

func (c *Client) GetIncidentSubStatus(id string) (*IncidentSubStatus, error) {
	req, err := rootlygo.NewGetIncidentSubStatusRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentSubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_sub_status: %s", err.Error())
	}

	return data.(*IncidentSubStatus), nil
}

func (c *Client) UpdateIncidentSubStatus(id string, incident_sub_status *IncidentSubStatus) (*IncidentSubStatus, error) {
	buffer, err := MarshalData(incident_sub_status)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_sub_status: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentSubStatusRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_sub_status: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentSubStatus))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_sub_status: %s", err.Error())
	}

	return data.(*IncidentSubStatus), nil
}

func (c *Client) DeleteIncidentSubStatus(id string) error {
	req, err := rootlygo.NewDeleteIncidentSubStatusRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_sub_status: %s", err.Error())
	}

	return nil
}