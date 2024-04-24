package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Schedule struct {
	ID string `jsonapi:"primary,schedules"`
	Name string `jsonapi:"attr,name,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
}

func (c *Client) ListSchedules(params *rootlygo.ListSchedulesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListSchedulesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	schedules, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Schedule)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return schedules, nil
}

func (c *Client) CreateSchedule(d *Schedule) (*Schedule, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule: %s", err.Error())
	}

	req, err := rootlygo.NewCreateScheduleRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create schedule: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Schedule))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule: %s", err.Error())
	}

	return data.(*Schedule), nil
}

func (c *Client) GetSchedule(id string) (*Schedule, error) {
	req, err := rootlygo.NewGetScheduleRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get schedule: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Schedule))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule: %s", err.Error())
	}

	return data.(*Schedule), nil
}

func (c *Client) UpdateSchedule(id string, schedule *Schedule) (*Schedule, error) {
	buffer, err := MarshalData(schedule)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateScheduleRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update schedule: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Schedule))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule: %s", err.Error())
	}

	return data.(*Schedule), nil
}

func (c *Client) DeleteSchedule(id string) error {
	req, err := rootlygo.NewDeleteScheduleRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete schedule: %s", err.Error())
	}

	return nil
}
