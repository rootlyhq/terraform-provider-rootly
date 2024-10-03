package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type OnCallShadow struct {
	ID string `jsonapi:"primary,on_call_shadows"`
	ScheduleId string `jsonapi:"attr,schedule_id,omitempty"`
  ShadowableType string `jsonapi:"attr,shadowable_type,omitempty"`
  ShadowableId string `jsonapi:"attr,shadowable_id,omitempty"`
  ShadowUserId int `jsonapi:"attr,shadow_user_id,omitempty"`
  StartsAt string `jsonapi:"attr,starts_at,omitempty"`
  EndsAt string `jsonapi:"attr,ends_at,omitempty"`
}

func (c *Client) ListOnCallShadows(id string, params *rootlygo.ListOnCallShadowsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListOnCallShadowsRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	on_call_shadows, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(OnCallShadow)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return on_call_shadows, nil
}

func (c *Client) CreateOnCallShadow(d *OnCallShadow) (*OnCallShadow, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling on_call_shadow: %s", err.Error())
	}

	req, err := rootlygo.NewCreateOnCallShadowRequestWithBody(c.Rootly.Server, d.ScheduleId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create on_call_shadow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OnCallShadow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling on_call_shadow: %s", err.Error())
	}

	return data.(*OnCallShadow), nil
}

func (c *Client) GetOnCallShadow(id string) (*OnCallShadow, error) {
	req, err := rootlygo.NewGetOnCallShadowRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get on_call_shadow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OnCallShadow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling on_call_shadow: %s", err.Error())
	}

	return data.(*OnCallShadow), nil
}

func (c *Client) UpdateOnCallShadow(id string, on_call_shadow *OnCallShadow) (*OnCallShadow, error) {
	buffer, err := MarshalData(on_call_shadow)
	if err != nil {
		return nil, errors.Errorf("Error marshaling on_call_shadow: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateOnCallShadowRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update on_call_shadow: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(OnCallShadow))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling on_call_shadow: %s", err.Error())
	}

	return data.(*OnCallShadow), nil
}

func (c *Client) DeleteOnCallShadow(id string) error {
	req, err := rootlygo.NewDeleteOnCallShadowRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete on_call_shadow: %s", err.Error())
	}

	return nil
}
