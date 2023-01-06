package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Secret struct {
	ID string `jsonapi:"primary,secrets"`
	Name string `jsonapi:"attr,name,omitempty"`
  Secret string `jsonapi:"attr,secret,omitempty"`
}

func (c *Client) ListSecrets(params *rootlygo.ListSecretsParams) ([]interface{}, error) {
	req, err := rootlygo.NewListSecretsRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	secrets, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Secret)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return secrets, nil
}

func (c *Client) CreateSecret(d *Secret) (*Secret, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling secret: %s", err.Error())
	}

	req, err := rootlygo.NewCreateSecretRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create secret: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling secret: %s", err.Error())
	}

	return data.(*Secret), nil
}

func (c *Client) GetSecret(id string) (*Secret, error) {
	req, err := rootlygo.NewGetSecretRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get secret: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling secret: %s", err.Error())
	}

	return data.(*Secret), nil
}

func (c *Client) UpdateSecret(id string, secret *Secret) (*Secret, error) {
	buffer, err := MarshalData(secret)
	if err != nil {
		return nil, errors.Errorf("Error marshaling secret: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateSecretRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update secret: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Secret))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling secret: %s", err.Error())
	}

	return data.(*Secret), nil
}

func (c *Client) DeleteSecret(id string) error {
	req, err := rootlygo.NewDeleteSecretRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete secret: %s", err.Error())
	}

	return nil
}
