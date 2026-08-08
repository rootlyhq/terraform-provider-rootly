// Hand-maintained; kept out of codegen via excluded.clients in tools/generate.js.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

// Captures only the relationship id; no `include` sideload required.
type onCallRoleRef struct {
	ID string `jsonapi:"primary,on_call_roles"`
}

type roleRef struct {
	ID string `jsonapi:"primary,roles"`
}

type User struct {
	ID               string         `jsonapi:"primary,users"`
	Email            string         `jsonapi:"attr,email,omitempty"`
	FirstName        string         `jsonapi:"attr,first_name,omitempty"`
	LastName         string         `jsonapi:"attr,last_name,omitempty"`
	FullName         string         `jsonapi:"attr,full_name,omitempty"`
	FullNameWithTeam string         `jsonapi:"attr,full_name_with_team,omitempty"`
	TimeZone         string         `jsonapi:"attr,time_zone,omitempty"`
	OnCallRole       *onCallRoleRef `jsonapi:"relation,on_call_role,omitempty"`
	Role             *roleRef       `jsonapi:"relation,role,omitempty"`
}

func (u *User) OnCallRoleId() string {
	if u.OnCallRole == nil {
		return ""
	}
	return u.OnCallRole.ID
}

func (u *User) RoleId() string {
	if u.Role == nil {
		return ""
	}
	return u.Role.ID
}

func (c *Client) ListUsers(params *rootlygo.ListUsersParams) ([]interface{}, error) {
	if params == nil {
		params = &rootlygo.ListUsersParams{}
	}
	// Always fetch the role relationships so callers get OnCallRoleId/RoleId.
	if params.Include == nil {
		include := rootlygo.ListUsersParamsInclude("on_call_role,role")
		params.Include = &include
	}

	req, err := rootlygo.NewListUsersRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	users, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(User)))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling: %w", err)
	}

	return users, nil
}

func (c *Client) GetUser(id string) (*User, error) {
	include := rootlygo.GetUserParamsInclude("on_call_role,role")
	req, err := rootlygo.NewGetUserRequest(c.Rootly.Server, id, &rootlygo.GetUserParams{Include: &include})
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get user: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(User))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling user: %w", err)
	}

	return data.(*User), nil
}

// Minimal payload: only on_call_role_id is sent so other user attributes are
// left untouched by the API.
type updateUserOnCallRoleBody struct {
	Data updateUserOnCallRoleData `json:"data"`
}

type updateUserOnCallRoleData struct {
	Type       string                         `json:"type"`
	Attributes updateUserOnCallRoleAttributes `json:"attributes"`
}

type updateUserOnCallRoleAttributes struct {
	OnCallRoleId *string `json:"on_call_role_id"`
}

// UpdateUserOnCallRole sets a user's on-call role, or clears it when
// onCallRoleId is empty (nil pointer marshals to JSON null).
func (c *Client) UpdateUserOnCallRole(id string, onCallRoleId string) (*User, error) {
	attrs := updateUserOnCallRoleAttributes{}
	if onCallRoleId != "" {
		attrs.OnCallRoleId = &onCallRoleId
	}

	body := updateUserOnCallRoleBody{
		Data: updateUserOnCallRoleData{
			Type:       "users",
			Attributes: attrs,
		},
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling user on_call_role update: %w", err)
	}

	req, err := rootlygo.NewUpdateUserRequestWithBody(c.Rootly.Server, id, c.ContentType, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update user on_call_role: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(User))
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling user: %w", err)
	}

	return data.(*User), nil
}
