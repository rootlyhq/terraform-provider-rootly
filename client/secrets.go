package client

import (
	"reflect"

	"fmt"

	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

type Secret struct {
	ID                    string `jsonapi:"primary,secrets"`
	Name                  string `jsonapi:"attr,name,omitempty"`
	Secret                string `jsonapi:"attr,secret,omitempty"`
	HashicorpVaultMount   string `jsonapi:"attr,hashicorp_vault_mount,omitempty"`
	HashicorpVaultPath    string `jsonapi:"attr,hashicorp_vault_path,omitempty"`
	HashicorpVaultVersion int    `jsonapi:"attr,hashicorp_vault_version,omitempty"`
}

func (c *Client) ListSecrets(params *rootlygo.ListSecretsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListSecretsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}

	secrets, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Secret)))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling: %w", err)
	}

	return secrets, nil
}

func (c *Client) CreateSecret(d *Secret) (*Secret, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling secret: %w", err)
	}

	req, err := rootlygo.NewCreateSecretRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request to create secret: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling secret: %w", err)
	}

	return data.(*Secret), nil
}

func (c *Client) GetSecret(id string) (*Secret, error) {
	req, err := rootlygo.NewGetSecretRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to get secret: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling secret: %w", err)
	}

	return data.(*Secret), nil
}

func (c *Client) UpdateSecret(id string, secret *Secret) (*Secret, error) {
	buffer, err := MarshalData(secret)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling secret: %w", err)
	}

	req, err := rootlygo.NewUpdateSecretRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error building request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to update secret: %w", err)
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling secret: %w", err)
	}

	return data.(*Secret), nil
}

func (c *Client) DeleteSecret(id string) error {
	req, err := rootlygo.NewDeleteSecretRequest(c.Rootly.Server, id)
	if err != nil {
		return fmt.Errorf("Error building request: %w", err)
	}

	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request to delete secret: %w", err)
	}

	return nil
}
