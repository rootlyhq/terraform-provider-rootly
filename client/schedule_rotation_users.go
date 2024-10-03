package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type ScheduleRotationUser struct {
	ID string `jsonapi:"primary,schedule_rotation_users"`
	ScheduleRotationId string `jsonapi:"attr,schedule_rotation_id,omitempty"`
  UserId int `jsonapi:"attr,user_id,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
}

func (c *Client) ListScheduleRotationUsers(id string, params *rootlygo.ListScheduleRotationUsersParams) ([]interface{}, error) {
	req, err := rootlygo.NewListScheduleRotationUsersRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	schedule_rotation_users, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(ScheduleRotationUser)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return schedule_rotation_users, nil
}

func (c *Client) CreateScheduleRotationUser(d *ScheduleRotationUser) (*ScheduleRotationUser, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_user: %s", err.Error())
	}

	req, err := rootlygo.NewCreateScheduleRotationUserRequestWithBody(c.Rootly.Server, d.ScheduleRotationId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create schedule_rotation_user: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationUser))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_user: %s", err.Error())
	}

	return data.(*ScheduleRotationUser), nil
}

func (c *Client) GetScheduleRotationUser(id string) (*ScheduleRotationUser, error) {
	req, err := rootlygo.NewGetScheduleRotationUserRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get schedule_rotation_user: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationUser))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_user: %s", err.Error())
	}

	return data.(*ScheduleRotationUser), nil
}

func (c *Client) UpdateScheduleRotationUser(id string, schedule_rotation_user *ScheduleRotationUser) (*ScheduleRotationUser, error) {
	buffer, err := MarshalData(schedule_rotation_user)
	if err != nil {
		return nil, errors.Errorf("Error marshaling schedule_rotation_user: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateScheduleRotationUserRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update schedule_rotation_user: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(ScheduleRotationUser))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling schedule_rotation_user: %s", err.Error())
	}

	return data.(*ScheduleRotationUser), nil
}

func (c *Client) DeleteScheduleRotationUser(id string) error {
	req, err := rootlygo.NewDeleteScheduleRotationUserRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete schedule_rotation_user: %s", err.Error())
	}

	return nil
}
