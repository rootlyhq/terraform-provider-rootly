package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type ScheduleRotationActiveTime struct {
	ID string `jsonapi:"primary,schedule_rotation_active_times"`
	ScheduleRotationId string `jsonapi:"attr,schedule_rotation_id,omitempty"`
  StartTime string `jsonapi:"attr,start_time,omitempty"`
  EndTime string `jsonapi:"attr,end_time,omitempty"`
}

func (c *Client) ListScheduleRotationActiveTimes(id string, params *rootlygo.ListScheduleRotationActiveTimesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListScheduleRotationActiveTimesRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	schedule_rotation_active_times, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(ScheduleRotationActiveTime)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return schedule_rotation_active_times, nil
}

func (c *Client) CreateScheduleRotationActiveTime(d *ScheduleRotationActiveTime) (*ScheduleRotationActiveTime, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_active_time: %s", err.Error())
	}

	req, err := rootlygo.NewCreateScheduleRotationActiveTimeRequestWithBody(c.Rootly.Server, d.ScheduleRotationId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create schedule_rotation_active_time: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveTime))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_time: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveTime), nil
}

func (c *Client) GetScheduleRotationActiveTime(id string) (*ScheduleRotationActiveTime, error) {
	req, err := rootlygo.NewGetScheduleRotationActiveTimeRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get schedule_rotation_active_time: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveTime))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_time: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveTime), nil
}

func (c *Client) UpdateScheduleRotationActiveTime(id string, schedule_rotation_active_time *ScheduleRotationActiveTime) (*ScheduleRotationActiveTime, error) {
	buffer, err := MarshalData(schedule_rotation_active_time)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_active_time: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateScheduleRotationActiveTimeRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update schedule_rotation_active_time: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveTime))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_time: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveTime), nil
}

func (c *Client) DeleteScheduleRotationActiveTime(id string) error {
	req, err := rootlygo.NewDeleteScheduleRotationActiveTimeRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete schedule_rotation_active_time: %s", err.Error())
	}

	return nil
}
