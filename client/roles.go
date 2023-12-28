package client

import (
	"reflect"
	
	"github.com/pkg/errors"
	"github.com/google/jsonapi"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

type Role struct {
	ID string `jsonapi:"primary,roles"`
	Name string `jsonapi:"attr,name,omitempty"`
  Slug string `jsonapi:"attr,slug,omitempty"`
  IncidentPermissionSetId string `jsonapi:"attr,incident_permission_set_id,omitempty"`
  IsDeletable *bool `jsonapi:"attr,is_deletable,omitempty"`
  IsEditable *bool `jsonapi:"attr,is_editable,omitempty"`
  ApiKeysPermissions []interface{} `jsonapi:"attr,api_keys_permissions,omitempty"`
  AuditsPermissions []interface{} `jsonapi:"attr,audits_permissions,omitempty"`
  BillingPermissions []interface{} `jsonapi:"attr,billing_permissions,omitempty"`
  EnvironmentsPermissions []interface{} `jsonapi:"attr,environments_permissions,omitempty"`
  FormFieldsPermissions []interface{} `jsonapi:"attr,form_fields_permissions,omitempty"`
  FunctionalitiesPermissions []interface{} `jsonapi:"attr,functionalities_permissions,omitempty"`
  GroupsPermissions []interface{} `jsonapi:"attr,groups_permissions,omitempty"`
  IncidentCausesPermissions []interface{} `jsonapi:"attr,incident_causes_permissions,omitempty"`
  IncidentFeedbacksPermissions []interface{} `jsonapi:"attr,incident_feedbacks_permissions,omitempty"`
  IncidentPostMortemsPermissions []interface{} `jsonapi:"attr,incident_post_mortems_permissions,omitempty"`
  IncidentRolesPermissions []interface{} `jsonapi:"attr,incident_roles_permissions,omitempty"`
  IncidentTypesPermissions []interface{} `jsonapi:"attr,incident_types_permissions,omitempty"`
  IncidentsPermissions []interface{} `jsonapi:"attr,incidents_permissions,omitempty"`
  InvitationsPermissions []interface{} `jsonapi:"attr,invitations_permissions,omitempty"`
  PlaybooksPermissions []interface{} `jsonapi:"attr,playbooks_permissions,omitempty"`
  PrivateIncidentsPermissions []interface{} `jsonapi:"attr,private_incidents_permissions,omitempty"`
  RetrospectivePermissions []interface{} `jsonapi:"attr,retrospective_permissions,omitempty"`
  RolesPermissions []interface{} `jsonapi:"attr,roles_permissions,omitempty"`
  SecretsPermissions []interface{} `jsonapi:"attr,secrets_permissions,omitempty"`
  ServicesPermissions []interface{} `jsonapi:"attr,services_permissions,omitempty"`
  SeveritiesPermissions []interface{} `jsonapi:"attr,severities_permissions,omitempty"`
  StatusPagesPermissions []interface{} `jsonapi:"attr,status_pages_permissions,omitempty"`
  WebhooksPermissions []interface{} `jsonapi:"attr,webhooks_permissions,omitempty"`
  WorkflowsPermissions []interface{} `jsonapi:"attr,workflows_permissions,omitempty"`
}

func (c *Client) ListRoles(params *rootlygo.ListRolesParams) ([]interface{}, error) {
	req, err := rootlygo.NewListRolesRequest(c.Rootly.Server, params)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request: %s", err.Error())
	}

	roles, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Role)))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling: %s", err.Error())
	}

	return roles, nil
}

func (c *Client) CreateRole(d *Role) (*Role, error) {
	buffer, err := MarshalData(d)
	if err != nil {
		return nil, errors.Errorf("Error marshaling role: %s", err.Error())
	}

	req, err := rootlygo.NewCreateRoleRequestWithBody(c.Rootly.Server, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to perform request to create role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Role))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling role: %s", err.Error())
	}

	return data.(*Role), nil
}

func (c *Client) GetRole(id string) (*Role, error) {
	req, err := rootlygo.NewGetRoleRequest(c.Rootly.Server, id)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to get role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Role))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling role: %s", err.Error())
	}

	return data.(*Role), nil
}

func (c *Client) UpdateRole(id string, role *Role) (*Role, error) {
	buffer, err := MarshalData(role)
	if err != nil {
		return nil, errors.Errorf("Error marshaling role: %s", err.Error())
	}

	req, err := rootlygo.NewUpdateRoleRequestWithBody(c.Rootly.Server, id, c.ContentType, buffer)
	if err != nil {
		return nil, errors.Errorf("Error building request: %s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Errorf("Failed to make request to update role: %s", err.Error())
	}

	data, err := UnmarshalData(resp.Body, new(Role))
	if err != nil {
		return nil, errors.Errorf("Error unmarshaling role: %s", err.Error())
	}

	return data.(*Role), nil
}

func (c *Client) DeleteRole(id string) error {
	req, err := rootlygo.NewDeleteRoleRequest(c.Rootly.Server, id)
	if err != nil {
		return errors.Errorf("Error building request: %s", err.Error())
	}

	_, err = c.Do(req)
	if err != nil {
		return errors.Errorf("Failed to make request to delete role: %s", err.Error())
	}

	return nil
}
