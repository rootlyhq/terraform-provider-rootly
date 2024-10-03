package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type ScheduleRotationActiveDay struct {
	ID string `jsonapi:"primary,schedule_rotation_active_days"`
	ScheduleRotationId string `jsonapi:"attr,schedule_rotation_id,omitempty"`
  DayName string `jsonapi:"attr,day_name,omitempty"`
  ActiveTimeAttributes []interface{} `jsonapi:"attr,active_time_attributes,omitempty"`
}

func (c *Client) ListScheduleRotationActiveDays(id string, params *rootlygo.ListScheduleRotationActiveDaysParams) ([]interface{}, error) {
	req, err := rootlygo.NewListScheduleRotationActiveDaysRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	schedule_rotation_active_days, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(ScheduleRotationActiveDay)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return schedule_rotation_active_days, nil
}

func (c *Client) CreateScheduleRotationActiveDay(d *ScheduleRotationActiveDay) (*ScheduleRotationActiveDay, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_active_day: %s", err.Error())
	}

	req, err := rootlygo.NewCreateScheduleRotationActiveDayRequestWithBody(c.Rootly.Server, d.ScheduleRotationId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create schedule_rotation_active_day: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveDay))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_day: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveDay), nil
}

func (c *Client) GetScheduleRotationActiveDay(id string) (*ScheduleRotationActiveDay, error) {
	req, err := rootlygo.NewGetScheduleRotationActiveDayRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get schedule_rotation_active_day: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveDay))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_day: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveDay), nil
}

func (c *Client) UpdateScheduleRotationActiveDay(id string, schedule_rotation_active_day *ScheduleRotationActiveDay) (*ScheduleRotationActiveDay, error) {
	buffer, err := MarshalData(schedule_rotation_active_day)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_active_day: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateScheduleRotationActiveDayRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update schedule_rotation_active_day: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationActiveDay))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_active_day: %s", err.Error())
	}

	return data.(*ScheduleRotationActiveDay), nil
}

func (c *Client) DeleteScheduleRotationActiveDay(id string) error {
	req, err := rootlygo.NewDeleteScheduleRotationActiveDayRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete schedule_rotation_active_day: %s", err.Error())
	}

	return nil
}
