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

func resourceIncidentPermissionSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentPermissionSetCreate,
		ReadContext:   resourceIncidentPermissionSetRead,
		UpdateContext: resourceIncidentPermissionSetUpdate,
		DeleteContext: resourceIncidentPermissionSetDelete,
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
				Description: "The incident permission set name.",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The incident permission set slug.",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The incident permission set description.",
			},

			"private_incident_permissions": &schema.Schema{
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

			"public_incident_permissions": &schema.Schema{
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

func resourceIncidentPermissionSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentPermissionSet"))

	s := &client.IncidentPermissionSet{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("private_incident_permissions"); ok {
		s.PrivateIncidentPermissions = value.([]interface{})
	}
	if value, ok := d.GetOkExists("public_incident_permissions"); ok {
		s.PublicIncidentPermissions = value.([]interface{})
	}

	res, err := c.CreateIncidentPermissionSet(s)
	if err != nil {
		return diag.Errorf("Error creating incident_permission_set: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_permission_set resource: %s", d.Id()))

	return resourceIncidentPermissionSetRead(ctx, d, meta)
}

func resourceIncidentPermissionSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentPermissionSet: %s", d.Id()))

	item, err := c.GetIncidentPermissionSet(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSet (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_permission_set: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("private_incident_permissions", item.PrivateIncidentPermissions)
	d.Set("public_incident_permissions", item.PublicIncidentPermissions)

	return nil
}

func resourceIncidentPermissionSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentPermissionSet: %s", d.Id()))

	s := &client.IncidentPermissionSet{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("private_incident_permissions") {
		s.PrivateIncidentPermissions = d.Get("private_incident_permissions").([]interface{})
	}
	if d.HasChange("public_incident_permissions") {
		s.PublicIncidentPermissions = d.Get("public_incident_permissions").([]interface{})
	}

	_, err := c.UpdateIncidentPermissionSet(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_permission_set: %s", err.Error())
	}

	return resourceIncidentPermissionSetRead(ctx, d, meta)
}

func resourceIncidentPermissionSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentPermissionSet: %s", d.Id()))

	err := c.DeleteIncidentPermissionSet(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSet (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_permission_set: %s", err.Error())
	}

	d.SetId("")

	return nil
}
