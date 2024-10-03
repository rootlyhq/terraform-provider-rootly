package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Authorization struct {
	ID string `jsonapi:"primary,authorizations"`
	AuthorizableId string `jsonapi:"attr,authorizable_id,omitempty"`
  AuthorizableType string `jsonapi:"attr,authorizable_type,omitempty"`
  GranteeId string `jsonapi:"attr,grantee_id,omitempty"`
  GranteeType string `jsonapi:"attr,grantee_type,omitempty"`
  Permissions []interface{} `jsonapi:"attr,permissions,omitempty"`
}

func (c *Client) ListAuthorizations(params *rootlygo.ListAuthorizationsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListAuthorizationsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	authorizations, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Authorization)))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return authorizations, nil
}

func (c *Client) CreateAuthorization(d *Authorization) (*Authorization, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling authorization: %s", err.Error())
	}

	req, err := rootlygo.NewCreateAuthorizationRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create authorization: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Authorization))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling authorization: %s", err.Error())
	}

	return data.(*Authorization), nil
}

func (c *Client) GetAuthorization(id string) (*Authorization, error) {
	req, err := rootlygo.NewGetAuthorizationRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get authorization: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Authorization))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling authorization: %s", err.Error())
	}

	return data.(*Authorization), nil
}

func (c *Client) UpdateAuthorization(id string, authorization *Authorization) (*Authorization, error) {
	buffer, err := MarshalData(authorization)
	if err != nil {
		return nil, errors.Errorf("Error marshaling authorization: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateAuthorizationRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update authorization: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Authorization))
	resp.Body.Close()
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling authorization: %s", err.Error())
	}

	return data.(*Authorization), nil
}

func (c *Client) DeleteAuthorization(id string) error {
	req, err := rootlygo.NewDeleteAuthorizationRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete authorization: %s", err.Error())
	}

	return nil
}
