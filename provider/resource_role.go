package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The role name.",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The role slug.",
			},

			"incident_permission_set_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Associated incident permissions set.",
			},

			"is_deletable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether the role can be deleted.. Value must be one of true or false",
				Deprecated:  "This resource is now ignored by the API, and will be removed in the next release.",
			},

			"is_editable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether the role can be edited.. Value must be one of true or false",
				Deprecated:  "This resource is now ignored by the API, and will be removed in the next release.",
			},

			"alerts_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`.",
			},

			"pulses_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `update`, `read`.",
			},

			"api_keys_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"audits_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"billing_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"environments_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"form_fields_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"functionalities_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"groups_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"incident_causes_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"incident_feedbacks_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"incident_roles_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"incident_types_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"incidents_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"invitations_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"playbooks_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"private_incidents_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"retrospective_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"roles_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"secrets_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"services_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"severities_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"status_pages_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"webhooks_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},

			"workflows_permissions": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `create`, `read`, `update`, `delete`.",
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Role"))

	s := &client.Role{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("incident_permission_set_id"); ok {
		s.IncidentPermissionSetId = value.(string)
	}
	if value, ok := d.GetOkExists("alerts_permissions"); ok {
		s.AlertsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("pulses_permissions"); ok {
		s.PulsesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("api_keys_permissions"); ok {
		s.ApiKeysPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("audits_permissions"); ok {
		s.AuditsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("billing_permissions"); ok {
		s.BillingPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("environments_permissions"); ok {
		s.EnvironmentsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("form_fields_permissions"); ok {
		s.FormFieldsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("functionalities_permissions"); ok {
		s.FunctionalitiesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("groups_permissions"); ok {
		s.GroupsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_causes_permissions"); ok {
		s.IncidentCausesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_feedbacks_permissions"); ok {
		s.IncidentFeedbacksPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_roles_permissions"); ok {
		s.IncidentRolesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_types_permissions"); ok {
		s.IncidentTypesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incidents_permissions"); ok {
		s.IncidentsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("invitations_permissions"); ok {
		s.InvitationsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("playbooks_permissions"); ok {
		s.PlaybooksPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("private_incidents_permissions"); ok {
		s.PrivateIncidentsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("retrospective_permissions"); ok {
		s.RetrospectivePermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("roles_permissions"); ok {
		s.RolesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("secrets_permissions"); ok {
		s.SecretsPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("services_permissions"); ok {
		s.ServicesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("severities_permissions"); ok {
		s.SeveritiesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("status_pages_permissions"); ok {
		s.StatusPagesPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("webhooks_permissions"); ok {
		s.WebhooksPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("workflows_permissions"); ok {
		s.WorkflowsPermissions = value.([]interface{})
	}

	res, err := c.CreateRole(s)
	if err != nil {
		return diag.Errorf("Error creating role: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a role resource: %s", d.Id()))

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Role: %s", d.Id()))

	item, err := c.GetRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Role (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading role: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("incident_permission_set_id", item.IncidentPermissionSetId)
	d.Set("alerts_permissions", item.AlertsPermissions)
	d.Set("pulses_permissions", item.PulsesPermissions)
	d.Set("api_keys_permissions", item.ApiKeysPermissions)
	d.Set("audits_permissions", item.AuditsPermissions)
	d.Set("billing_permissions", item.BillingPermissions)
	d.Set("environments_permissions", item.EnvironmentsPermissions)
	d.Set("form_fields_permissions", item.FormFieldsPermissions)
	d.Set("functionalities_permissions", item.FunctionalitiesPermissions)
	d.Set("groups_permissions", item.GroupsPermissions)
	d.Set("incident_causes_permissions", item.IncidentCausesPermissions)
	d.Set("incident_feedbacks_permissions", item.IncidentFeedbacksPermissions)
	d.Set("incident_roles_permissions", item.IncidentRolesPermissions)
	d.Set("incident_types_permissions", item.IncidentTypesPermissions)
	d.Set("incidents_permissions", item.IncidentsPermissions)
	d.Set("invitations_permissions", item.InvitationsPermissions)
	d.Set("playbooks_permissions", item.PlaybooksPermissions)
	d.Set("private_incidents_permissions", item.PrivateIncidentsPermissions)
	d.Set("retrospective_permissions", item.RetrospectivePermissions)
	d.Set("roles_permissions", item.RolesPermissions)
	d.Set("secrets_permissions", item.SecretsPermissions)
	d.Set("services_permissions", item.ServicesPermissions)
	d.Set("severities_permissions", item.SeveritiesPermissions)
	d.Set("status_pages_permissions", item.StatusPagesPermissions)
	d.Set("webhooks_permissions", item.WebhooksPermissions)
	d.Set("workflows_permissions", item.WorkflowsPermissions)

	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Role: %s", d.Id()))

	s := &client.Role{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("incident_permission_set_id") {
		s.IncidentPermissionSetId = d.Get("incident_permission_set_id").(string)
	}
	if d.HasChange("alerts_permissions") {
		s.AlertsPermissions = d.Get("alerts_permissions").([]interface{})
	}
	if d.HasChange("pulses_permissions") {
		s.PulsesPermissions = d.Get("pulses_permissions").([]interface{})
	}
	if d.HasChange("api_keys_permissions") {
		s.ApiKeysPermissions = d.Get("api_keys_permissions").([]interface{})
	}
	if d.HasChange("audits_permissions") {
		s.AuditsPermissions = d.Get("audits_permissions").([]interface{})
	}
	if d.HasChange("billing_permissions") {
		s.BillingPermissions = d.Get("billing_permissions").([]interface{})
	}
	if d.HasChange("environments_permissions") {
		s.EnvironmentsPermissions = d.Get("environments_permissions").([]interface{})
	}
	if d.HasChange("form_fields_permissions") {
		s.FormFieldsPermissions = d.Get("form_fields_permissions").([]interface{})
	}
	if d.HasChange("functionalities_permissions") {
		s.FunctionalitiesPermissions = d.Get("functionalities_permissions").([]interface{})
	}
	if d.HasChange("groups_permissions") {
		s.GroupsPermissions = d.Get("groups_permissions").([]interface{})
	}
	if d.HasChange("incident_causes_permissions") {
		s.IncidentCausesPermissions = d.Get("incident_causes_permissions").([]interface{})
	}
	if d.HasChange("incident_feedbacks_permissions") {
		s.IncidentFeedbacksPermissions = d.Get("incident_feedbacks_permissions").([]interface{})
	}
	if d.HasChange("incident_roles_permissions") {
		s.IncidentRolesPermissions = d.Get("incident_roles_permissions").([]interface{})
	}
	if d.HasChange("incident_types_permissions") {
		s.IncidentTypesPermissions = d.Get("incident_types_permissions").([]interface{})
	}
	if d.HasChange("incidents_permissions") {
		s.IncidentsPermissions = d.Get("incidents_permissions").([]interface{})
	}
	if d.HasChange("invitations_permissions") {
		s.InvitationsPermissions = d.Get("invitations_permissions").([]interface{})
	}
	if d.HasChange("playbooks_permissions") {
		s.PlaybooksPermissions = d.Get("playbooks_permissions").([]interface{})
	}
	if d.HasChange("private_incidents_permissions") {
		s.PrivateIncidentsPermissions = d.Get("private_incidents_permissions").([]interface{})
	}
	if d.HasChange("retrospective_permissions") {
		s.RetrospectivePermissions = d.Get("retrospective_permissions").([]interface{})
	}
	if d.HasChange("roles_permissions") {
		s.RolesPermissions = d.Get("roles_permissions").([]interface{})
	}
	if d.HasChange("secrets_permissions") {
		s.SecretsPermissions = d.Get("secrets_permissions").([]interface{})
	}
	if d.HasChange("services_permissions") {
		s.ServicesPermissions = d.Get("services_permissions").([]interface{})
	}
	if d.HasChange("severities_permissions") {
		s.SeveritiesPermissions = d.Get("severities_permissions").([]interface{})
	}
	if d.HasChange("status_pages_permissions") {
		s.StatusPagesPermissions = d.Get("status_pages_permissions").([]interface{})
	}
	if d.HasChange("webhooks_permissions") {
		s.WebhooksPermissions = d.Get("webhooks_permissions").([]interface{})
	}
	if d.HasChange("workflows_permissions") {
		s.WorkflowsPermissions = d.Get("workflows_permissions").([]interface{})
	}

	_, err := c.UpdateRole(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating role: %s", err.Error())
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Role: %s", d.Id()))

	err := c.DeleteRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Role (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting role: %s", err.Error())
	}

	d.SetId("")

	return nil
}
