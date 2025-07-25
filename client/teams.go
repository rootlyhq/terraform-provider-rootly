// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package client

import (
    "fmt"
	"reflect"
	
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Team struct {
	ID string `jsonapi:"primary,groups"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  Description string `jsonapi:"attr,description,omitempty"`
  NotifyEmails []interface{} `jsonapi:"attr,notify_emails,omitempty"`
  Color string `jsonapi:"attr,color,omitempty"`
  Position int `jsonapi:"attr,position,omitempty"`
  BackstageId string `jsonapi:"attr,backstage_id,omitempty"`
  ExternalId string `jsonapi:"attr,external_id,omitempty"`
  PagerdutyId string `jsonapi:"attr,pagerduty_id,omitempty"`
  PagerdutyServiceId string `jsonapi:"attr,pagerduty_service_id,omitempty"`
  OpsgenieId string `jsonapi:"attr,opsgenie_id,omitempty"`
  VictorOpsId string `jsonapi:"attr,victor_ops_id,omitempty"`
  PagertreeId string `jsonapi:"attr,pagertree_id,omitempty"`
  CortexId string `jsonapi:"attr,cortex_id,omitempty"`
  ServiceNowCiSysId string `jsonapi:"attr,service_now_ci_sys_id,omitempty"`
  UserIds []interface{} `jsonapi:"attr,user_ids,omitempty"`
  AdminIds []interface{} `jsonapi:"attr,admin_ids,omitempty"`
  AlertsEmailEnabled *bool `jsonapi:"attr,alerts_email_enabled,omitempty"`
  AlertsEmailAddress string `jsonapi:"attr,alerts_email_address,omitempty"`
  AlertUrgencyId string `jsonapi:"attr,alert_urgency_id,omitempty"`
  SlackChannels []interface{} `jsonapi:"attr,slack_channels,omitempty"`
  SlackAliases []interface{} `jsonapi:"attr,slack_aliases,omitempty"`
}

func (c *Client) ListTeams(params *rootlygo.ListTeamsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListTeamsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	teams, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Team)))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return teams, nil
}

func (c *Client) CreateTeam(d *Team) (*Team, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling team: %w", err)
	}

	req, err := rootlygo.NewCreateTeamRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create team: %s", err)
	}

	data, err := UnmarshalData(resp.Body, new(Team))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling team: %w", err)
	}

	return data.(*Team), nil
}

func (c *Client) GetTeam(id string) (*Team, error) {
	req, err := rootlygo.NewGetTeamRequest(c.Rootly.Server, id, nil)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get team: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Team))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling team: %w", err)
	}

	return data.(*Team), nil
}

func (c *Client) UpdateTeam(id string, team *Team) (*Team, error) {
	buffer, err := MarshalData(team)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling team: %w", err)
	}

	req, err := rootlygo.NewUpdateTeamRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update team: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Team))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling team: %w", err)
	}

	return data.(*Team), nil
}

func (c *Client) DeleteTeam(id string) error {
	req, err := rootlygo.NewDeleteTeamRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete team: %w", err)
	}

	return nil
}
