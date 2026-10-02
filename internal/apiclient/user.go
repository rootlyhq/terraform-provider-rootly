package apiclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/rootlyhq/jsonapi"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type User struct {
	ID     string                       `jsonapi:"primary,users"`
	RoleId jsonapi.NullableAttr[string] `jsonapi:"attr,role_id"`
	Role   *Role                        `jsonapi:"relation,role,omitempty"`
}

type Role struct {
	ID string `jsonapi:"primary,roles"`
}

func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	include := rootly.GetUserParamsInclude("role")
	httpResp, err := c.ClientWithResponses.GetUserWithResponse(ctx, id, &rootly.GetUserParams{Include: &include})
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		return nil, client.NotFoundError{}
	}

	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode())
	}

	body := httpResp.Body

	var resp User
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateUser(ctx context.Context, req User) (*User, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.UpdateUserWithBodyWithResponse(ctx, req.ID, "application/vnd.api+json", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		return nil, client.NotFoundError{}
	}

	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode())
	}

	body := httpResp.Body

	var resp User
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
