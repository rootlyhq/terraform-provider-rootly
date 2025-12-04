package client

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type CommunicationsGroup struct {
	ID                                string        `jsonapi:"primary,communications_groups"`
	Name                              string        `jsonapi:"attr,name,omitempty"`
	Slug                              string        `jsonapi:"attr,slug,omitempty"`
	Description                       string        `jsonapi:"attr,description,omitempty"`
	CommunicationTypeId               string        `jsonapi:"attr,communication_type_id,omitempty"`
	IsPrivate                         *bool         `jsonapi:"attr,is_private,omitempty"`
	ConditionType                     string        `jsonapi:"attr,condition_type,omitempty"`
	SmsChannel                        *bool         `jsonapi:"attr,sms_channel,omitempty"`
	EmailChannel                      *bool         `jsonapi:"attr,email_channel,omitempty"`
	CommunicationGroupConditions      []interface{} `jsonapi:"attr,communication_group_conditions,omitempty"`
	CommunicationGroupMembers         []interface{} `jsonapi:"attr,communication_group_members,omitempty"`
	CommunicationExternalGroupMembers []interface{} `jsonapi:"attr,communication_external_group_members,omitempty"`
}

func (c *Client) ListCommunicationsGroups(params *rootlygo.ListCommunicationsGroupsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListCommunicationsGroupsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	communications_groups, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(CommunicationsGroup)))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return communications_groups, nil
}

func (c *Client) CreateCommunicationsGroup(d *CommunicationsGroup) (*CommunicationsGroup, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling communications_group: %w", err)
	}

	req, err := rootlygo.NewCreateCommunicationsGroupRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create communications_group: %s", err)
	}

	data, err := UnmarshalData(resp.Body, new(CommunicationsGroup))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling communications_group: %w", err)
	}

	return data.(*CommunicationsGroup), nil
}

func (c *Client) GetCommunicationsGroup(id string) (*CommunicationsGroup, error) {
	req, err := rootlygo.NewGetCommunicationsGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get communications_group: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CommunicationsGroup))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling communications_group: %w", err)
	}

	return data.(*CommunicationsGroup), nil
}

func (c *Client) UpdateCommunicationsGroup(id string, communications_group *CommunicationsGroup) (*CommunicationsGroup, error) {
	buffer, err := MarshalData(communications_group)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling communications_group: %w", err)
	}

	req, err := rootlygo.NewUpdateCommunicationsGroupRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update communications_group: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(CommunicationsGroup))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling communications_group: %w", err)
	}

	return data.(*CommunicationsGroup), nil
}

func (c *Client) DeleteCommunicationsGroup(id string) error {
	req, err := rootlygo.NewDeleteCommunicationsGroupRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete communications_group: %w", err)
	}

	return nil
}
