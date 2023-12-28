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

func resourceIncidentPermissionSetResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentPermissionSetResourceCreate,
		ReadContext:   resourceIncidentPermissionSetResourceRead,
		UpdateContext: resourceIncidentPermissionSetResourceUpdate,
		DeleteContext: resourceIncidentPermissionSetResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"incident_permission_set_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},

			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "severities",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `severities`, `incident_types`, `statuses`.",
			},

			"private": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Value must be one of true or false",
			},

			"resource_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
		},
	}
}

func resourceIncidentPermissionSetResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentPermissionSetResource"))

	s := &client.IncidentPermissionSetResource{}

	if value, ok := d.GetOkExists("incident_permission_set_id"); ok {
		s.IncidentPermissionSetId = value.(string)
	}
	if value, ok := d.GetOkExists("kind"); ok {
		s.Kind = value.(string)
	}
	if value, ok := d.GetOkExists("private"); ok {
		s.Private = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("resource_id"); ok {
		s.ResourceId = value.(string)
	}
	if value, ok := d.GetOkExists("resource_type"); ok {
		s.ResourceType = value.(string)
	}

	res, err := c.CreateIncidentPermissionSetResource(s)
	if err != nil {
		return diag.Errorf("Error creating incident_permission_set_resource: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_permission_set_resource resource: %s", d.Id()))

	return resourceIncidentPermissionSetResourceRead(ctx, d, meta)
}

func resourceIncidentPermissionSetResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentPermissionSetResource: %s", d.Id()))

	item, err := c.GetIncidentPermissionSetResource(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSetResource (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_permission_set_resource: %s", d.Id())
	}

	d.Set("incident_permission_set_id", item.IncidentPermissionSetId)
	d.Set("kind", item.Kind)
	d.Set("private", item.Private)
	d.Set("resource_id", item.ResourceId)
	d.Set("resource_type", item.ResourceType)

	return nil
}

func resourceIncidentPermissionSetResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentPermissionSetResource: %s", d.Id()))

	s := &client.IncidentPermissionSetResource{}

	if d.HasChange("incident_permission_set_id") {
		s.IncidentPermissionSetId = d.Get("incident_permission_set_id").(string)
	}
	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
	if d.HasChange("private") {
		s.Private = tools.Bool(d.Get("private").(bool))
	}
	if d.HasChange("resource_id") {
		s.ResourceId = d.Get("resource_id").(string)
	}
	if d.HasChange("resource_type") {
		s.ResourceType = d.Get("resource_type").(string)
	}

	_, err := c.UpdateIncidentPermissionSetResource(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_permission_set_resource: %s", err.Error())
	}

	return resourceIncidentPermissionSetResourceRead(ctx, d, meta)
}

func resourceIncidentPermissionSetResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentPermissionSetResource: %s", d.Id()))

	err := c.DeleteIncidentPermissionSetResource(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentPermissionSetResource (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_permission_set_resource: %s", err.Error())
	}

	d.SetId("")

	return nil
}
