package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type LiveCallRouter struct {
	ID string `jsonapi:"primary,live_call_routers"`
	Kind string `jsonapi:"attr,kind,omitempty"`
  Enabled *bool `jsonapi:"attr,enabled,omitempty"`
  Name string `jsonapi:"attr,name,omitempty"`
  CountryCode string `jsonapi:"attr,country_code,omitempty"`
  PhoneType string `jsonapi:"attr,phone_type,omitempty"`
  PhoneNumber string `jsonapi:"attr,phone_number,omitempty"`
  VoicemailGreeting string `jsonapi:"attr,voicemail_greeting,omitempty"`
  CallerGreeting string `jsonapi:"attr,caller_greeting,omitempty"`
  WaitingMusicUrl string `jsonapi:"attr,waiting_music_url,omitempty"`
  SentToVoicemailDelay int `jsonapi:"attr,sent_to_voicemail_delay,omitempty"`
  ShouldRedirectToVoicemailOnNoAnswer *bool `jsonapi:"attr,should_redirect_to_voicemail_on_no_answer,omitempty"`
  EscalationLevelDelayInSeconds int `jsonapi:"attr,escalation_level_delay_in_seconds,omitempty"`
  ShouldAutoResolveAlertOnCallEnd *bool `jsonapi:"attr,should_auto_resolve_alert_on_call_end,omitempty"`
  AlertUrgencyId string `jsonapi:"attr,alert_urgency_id,omitempty"`
  EscalationPolicyTriggerParams map[string]interface{} `jsonapi:"attr,escalation_policy_trigger_params,omitempty"`
}

func (c *Client) ListLiveCallRouters(params *rootlygo.ListLiveCallRoutersParams) ([]interface{}, error) {
	req, err := rootlygo.NewListLiveCallRoutersRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	live_call_routers, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(LiveCallRouter)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return live_call_routers, nil
}

func (c *Client) CreateLiveCallRouter(d *LiveCallRouter) (*LiveCallRouter, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling live_call_router: %s", err.Error())
	}

	req, err := rootlygo.NewCreateLiveCallRouterRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create live_call_router: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(LiveCallRouter))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling live_call_router: %s", err.Error())
	}

	return data.(*LiveCallRouter), nil
}

func (c *Client) GetLiveCallRouter(id string) (*LiveCallRouter, error) {
	req, err := rootlygo.NewGetLiveCallRouterRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get live_call_router: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(LiveCallRouter))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling live_call_router: %s", err.Error())
	}

	return data.(*LiveCallRouter), nil
}

func (c *Client) UpdateLiveCallRouter(id string, live_call_router *LiveCallRouter) (*LiveCallRouter, error) {
	buffer, err := MarshalData(live_call_router)
	if err != nil {
		return nil, errors.Errorf("Error marshaling live_call_router: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateLiveCallRouterRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update live_call_router: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(LiveCallRouter))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling live_call_router: %s", err.Error())
	}

	return data.(*LiveCallRouter), nil
}

func (c *Client) DeleteLiveCallRouter(id string) error {
	req, err := rootlygo.NewDeleteLiveCallRouterRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete live_call_router: %s", err.Error())
	}

	return nil
}
