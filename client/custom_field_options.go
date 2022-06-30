package client

import (
	"strconv"
	"github.com/pkg/errors"
	rootlygo "github.com/rootlyhq/rootly-go"
)

type CustomFieldOption struct {
	ID            string `jsonapi:"primary,custom_field_options"`
	CustomFieldId int    `jsonapi:"attr,custom_field_id,omitempty"`
	Value         string `jsonapi:"attr,value,omitempty"`
	Color         string `jsonapi:"attr,color,omitempty"`
}

func (c *Client) CreateCustomFieldOption(i *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom field option: %s", err.Error())
	}

	req, err := rootlygo.NewCreateCustomFieldOptionRequestWithBody(c.Rootly.Server, strconv.Itoa(i.CustomFieldId), c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create custom field option: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) GetCustomFieldOption(id string) (*CustomFieldOption, error) {
	req, err := rootlygo.NewGetCustomFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get custom field option: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) UpdateCustomFieldOption(id string, i *CustomFieldOption) (*CustomFieldOption, error) {
	buffer, err := MarshalData(i)
	if err != nil {
		return nil, errors.Errorf("Error marshaling custom field option: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateCustomFieldOptionRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update custom field option: %s", id)
	}

	data, err := UnmarshalData(resp.Body, new(CustomFieldOption))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling custom field option: %s", err.Error())
	}

	return data.(*CustomFieldOption), nil
}

func (c *Client) DeleteCustomFieldOption(id string) error {
	req, err := rootlygo.NewDeleteCustomFieldOptionRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete custom field option: %s", id)
	}

	return nil
}
