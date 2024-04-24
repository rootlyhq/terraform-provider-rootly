package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type OverrideShift struct {
	ID string `jsonapi:"primary,override_shifts"`
	ScheduleId string `jsonapi:"attr,schedule_id,omitempty"`
  RotationId string `jsonapi:"attr,rotation_id,omitempty"`
  UserId int `jsonapi:"attr,user_id,omitempty"`
  StartsAt string `jsonapi:"attr,starts_at,omitempty"`
  EndsAt string `jsonapi:"attr,ends_at,omitempty"`
  IsOverride *bool `jsonapi:"attr,is_override,omitempty"`
  ShiftOverride map[string]interface{} `jsonapi:"attr,shift_override,omitempty"`
}

func (c *Client) ListOverrideShifts(id string, params *rootlygo.ListOverrideShiftsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListOverrideShiftsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	override_shifts, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(OverrideShift)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return override_shifts, nil
}

func (c *Client) CreateOverrideShift(d *OverrideShift) (*OverrideShift, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling override_shift: %s", err.Error())
	}

	req, err := rootlygo.NewCreateOverrideShiftRequestWithBody(c.Rootly.Server, d.ScheduleId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create override_shift: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OverrideShift))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling override_shift: %s", err.Error())
	}

	return data.(*OverrideShift), nil
}

func (c *Client) GetOverrideShift(id string) (*OverrideShift, error) {
	req, err := rootlygo.NewGetOverrideShiftRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get override_shift: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OverrideShift))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling override_shift: %s", err.Error())
	}

	return data.(*OverrideShift), nil
}

func (c *Client) UpdateOverrideShift(id string, override_shift *OverrideShift) (*OverrideShift, error) {
	buffer, err := MarshalData(override_shift)
	if err != nil {
		return nil, errors.Errorf("Error marshaling override_shift: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateOverrideShiftRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update override_shift: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OverrideShift))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling override_shift: %s", err.Error())
	}

	return data.(*OverrideShift), nil
}

func (c *Client) DeleteOverrideShift(id string) error {
	req, err := rootlygo.NewDeleteOverrideShiftRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete override_shift: %s", err.Error())
	}

	return nil
}
