package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type IncidentPermissionSetResource struct {
	ID string `jsonapi:"primary,incident_permission_set_resources"`
	IncidentPermissionSetId string `jsonapi:"attr,incident_permission_set_id,omitempty"`
  Kind string `jsonapi:"attr,kind,omitempty"`
  Private *bool `jsonapi:"attr,private,omitempty"`
  ResourceId string `jsonapi:"attr,resource_id,omitempty"`
  ResourceType string `jsonapi:"attr,resource_type,omitempty"`
}

func (c *Client) ListIncidentPermissionSetResources(id string, params *rootlygo.ListIncidentPermissionSetResourcesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListIncidentPermissionSetResourcesRequest(c.Rootly.Server, id, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	incident_permission_set_resources, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(IncidentPermissionSetResource)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return incident_permission_set_resources, nil
}

func (c *Client) CreateIncidentPermissionSetResource(d *IncidentPermissionSetResource) (*IncidentPermissionSetResource, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set_resource: %s", err.Error())
	}

	req, err := rootlygo.NewCreateIncidentPermissionSetResourceRequestWithBody(c.Rootly.Server, d.IncidentPermissionSetId, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create incident_permission_set_resource: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetResource))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_resource: %s", err.Error())
	}

	return data.(*IncidentPermissionSetResource), nil
}

func (c *Client) GetIncidentPermissionSetResource(id string) (*IncidentPermissionSetResource, error) {
	req, err := rootlygo.NewGetIncidentPermissionSetResourceRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get incident_permission_set_resource: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetResource))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_resource: %s", err.Error())
	}

	return data.(*IncidentPermissionSetResource), nil
}

func (c *Client) UpdateIncidentPermissionSetResource(id string, incident_permission_set_resource *IncidentPermissionSetResource) (*IncidentPermissionSetResource, error) {
	buffer, err := MarshalData(incident_permission_set_resource)
	if err != nil {
		return nil, errors.Errorf("Error marshaling incident_permission_set_resource: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateIncidentPermissionSetResourceRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update incident_permission_set_resource: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(IncidentPermissionSetResource))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling incident_permission_set_resource: %s", err.Error())
	}

	return data.(*IncidentPermissionSetResource), nil
}

func (c *Client) DeleteIncidentPermissionSetResource(id string) error {
	req, err := rootlygo.NewDeleteIncidentPermissionSetResourceRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete incident_permission_set_resource: %s", err.Error())
	}

	return nil
}
