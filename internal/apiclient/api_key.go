package apiclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rootlyhq/jsonapi"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

type ApiKey struct {
	ID                string                       `jsonapi:"primary,api_keys"`
	Name              string                       `jsonapi:"attr,name"`
	Kind              string                       `jsonapi:"attr,kind"`
	ExpiresAt         jsonapi.NullableAttr[string] `jsonapi:"attr,expires_at"`
	GracePeriodEndsAt jsonapi.NullableAttr[string] `jsonapi:"attr,grace_period_ends_at"`
	Token             string                       `jsonapi:"attr,token"`
}

// RotateApiKey is the rotate request body. The ID is left empty so that it is omitted from the payload; the API
// key to rotate is passed in the URL instead.
type RotateApiKey struct {
	ID                 string                       `jsonapi:"primary,api_keys"`
	ExpiresAt          jsonapi.NullableAttr[string] `jsonapi:"attr,expires_at"`
	GracePeriodMinutes jsonapi.NullableAttr[int64]  `jsonapi:"attr,grace_period_minutes"`
}

func (c *Client) GetApiKey(ctx context.Context, id string) (*ApiKey, error) {
	httpResp, err := c.ClientWithResponses.GetApiKeyWithResponse(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		return nil, client.NotFoundError{}
	}

	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode())
	}

	var resp ApiKey
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(httpResp.Body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RotateApiKey(ctx context.Context, id string, req RotateApiKey) (*ApiKey, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayloadWithoutIncluded(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.RotateApiKeyWithBody(ctx, id, "application/vnd.api+json", &buf)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, client.NotFoundError{}
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", httpResp.StatusCode, body)
	}

	var resp ApiKey
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateApiKey(ctx context.Context, req ApiKey) (*ApiKey, error) {
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayloadWithoutIncluded(&buf, &req); err != nil {
		return nil, err
	}

	httpResp, err := c.ClientWithResponses.CreateApiKeyWithBodyWithResponse(ctx, "application/vnd.api+json", &buf)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() != http.StatusCreated && httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", httpResp.StatusCode(), httpResp.Body)
	}

	var resp ApiKey
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(httpResp.Body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeleteApiKey(ctx context.Context, id string) error {
	httpResp, err := c.ClientWithResponses.DeleteApiKeyWithResponse(ctx, id)
	if err != nil {
		return err
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		return client.NotFoundError{}
	}

	if httpResp.StatusCode() != http.StatusOK && httpResp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode())
	}

	return nil
}
