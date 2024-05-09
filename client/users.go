package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type User struct {
	ID string `jsonapi:"primary,users"`
	Email string `jsonapi:"attr,email,omitempty"`
  FullName string `jsonapi:"attr,full_name,omitempty"`
  FullNameWithTeam string `jsonapi:"attr,full_name_with_team,omitempty"`
}

func (c *Client) ListUsers(params *rootlygo.ListUsersParams) ([]interface{}, error) {
	req, err := rootlygo.NewListUsersRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	users, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(User)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return users, nil
}

