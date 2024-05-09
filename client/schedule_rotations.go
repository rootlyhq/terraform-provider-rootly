package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type ScheduleRotation struct {
	ID string `jsonapi:"primary,schedule_rotations"`
	ScheduleId string `jsonapi:"attr,schedule_id,omitempty"`
  Name string `jsonapi:"attr,name,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  ScheduleRotationableType string `jsonapi:"attr,schedule_rotationable_type,omitempty"`
  ActiveAllWeek *bool `jsonapi:"attr,active_all_week,omitempty"`
  ActiveDays []interface{} `jsonapi:"attr,active_days,omitempty"`
  ActiveTimeType string `jsonapi:"attr,active_time_type,omitempty"`
  ActiveTimeAttributes []interface{} `jsonapi:"attr,active_time_attributes,omitempty"`
  TimeZone string `jsonapi:"attr,time_zone,omitempty"`
  ScheduleRotationableAttributes map[string]interface{} `jsonapi:"attr,schedule_rotationable_attributes,omitempty"`
}

func (c *Client) ListScheduleRotations(id string, params *rootlygo.ListScheduleRotationsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListScheduleRotationsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	schedule_rotations, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(ScheduleRotation)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return schedule_rotations, nil
}

func (c *Client) CreateScheduleRotation(d *ScheduleRotation) (*ScheduleRotation, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation: %s", err.Error())
	}

	req, err := rootlygo.NewCreateScheduleRotationRequestWithBody(c.Rootly.Server, d.ScheduleId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create schedule_rotation: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotation))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation: %s", err.Error())
	}

	return data.(*ScheduleRotation), nil
}

func (c *Client) GetScheduleRotation(id string) (*ScheduleRotation, error) {
	req, err := rootlygo.NewGetScheduleRotationRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get schedule_rotation: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotation))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation: %s", err.Error())
	}

	return data.(*ScheduleRotation), nil
}

func (c *Client) UpdateScheduleRotation(id string, schedule_rotation *ScheduleRotation) (*ScheduleRotation, error) {
	buffer, err := MarshalData(schedule_rotation)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateScheduleRotationRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update schedule_rotation: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotation))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation: %s", err.Error())
	}

	return data.(*ScheduleRotation), nil
}

func (c *Client) DeleteScheduleRotation(id string) error {
	req, err := rootlygo.NewDeleteScheduleRotationRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete schedule_rotation: %s", err.Error())
	}

	return nil
}
