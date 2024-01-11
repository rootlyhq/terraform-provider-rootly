package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	"github.com/rootlyhq/terraform-provider-rootly/tools"
)

func resourceIncidentPermissionSetBoolean() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentPermissionSetBooleanCreate,
		ReadContext:   resourceIncidentPermissionSetBooleanRead,
		UpdateContext: resourceIncidentPermissionSetBooleanUpdate,
		DeleteContext: resourceIncidentPermissionSetBooleanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"incident_permission_set_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "",
			},

			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "publish_to_status_page",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `publish_to_status_page`, `assign_incident_roles`, `invite_subscribers`, `update_summary`, `update_timeline`, `trigger_workflows`, `modify_custom_fields`.",
			},

			"private": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Value must be one of true or false",
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceIncidentPermissionSetBooleanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentPermissionSetBoolean"))

	s := &client.IncidentPermissionSetBoolean{}

	if value, ok := d.GetOkExists("incident_permission_set_id"); ok {
		s.IncidentPermissionSetId = value.(string)
	}
	if value, ok := d.GetOkExists("kind"); ok {
		s.Kind = value.(string)
	}
	if value, ok := d.GetOkExists("private"); ok {
		s.Private = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("enabled"); ok {
		s.Enabled = tools.Bool(value.(bool))
	}

	res, err := c.CreateIncidentPermissionSetBoolean(s)
	if err != nil {
		return diag.Errorf("Error creating incident_permission_set_boolean: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_permission_set_boolean resource: %s", d.Id()))

	return resourceIncidentPermissionSetBooleanRead(ctx, d, meta)
}

func resourceIncidentPermissionSetBooleanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentPermissionSetBoolean: %s", d.Id()))

	item, err := c.GetIncidentPermissionSetBoolean(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSetBoolean (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_permission_set_boolean: %s", d.Id())
	}

	d.Set("incident_permission_set_id", item.IncidentPermissionSetId)
	d.Set("kind", item.Kind)
	d.Set("private", item.Private)
	d.Set("enabled", item.Enabled)

	return nil
}

func resourceIncidentPermissionSetBooleanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentPermissionSetBoolean: %s", d.Id()))

	s := &client.IncidentPermissionSetBoolean{}

	if d.HasChange("incident_permission_set_id") {
		s.IncidentPermissionSetId = d.Get("incident_permission_set_id").(string)
	}
	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
	if d.HasChange("private") {
		s.Private = tools.Bool(d.Get("private").(bool))
	}
	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}

	_, err := c.UpdateIncidentPermissionSetBoolean(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_permission_set_boolean: %s", err.Error())
	}

	return resourceIncidentPermissionSetBooleanRead(ctx, d, meta)
}

func resourceIncidentPermissionSetBooleanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentPermissionSetBoolean: %s", d.Id()))

	err := c.DeleteIncidentPermissionSetBoolean(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSetBoolean (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_permission_set_boolean: %s", err.Error())
	}

	d.SetId("")

	return nil
}
