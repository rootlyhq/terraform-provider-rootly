package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootly/terraform-provider-rootly/client"
	"github.com/rootly/terraform-provider-rootly/tools"
)

func resourceIncidentRole() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Incident Roles (e.g Commander, Ops Lead, Communication).",

		CreateContext: resourceIncidentRoleCreate,
		ReadContext:   resourceIncidentRoleRead,
		UpdateContext: resourceIncidentRoleUpdate,
		DeleteContext: resourceIncidentRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the incident role",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the incident role",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"summary": {
				Description: "The summary of the incident role",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the incident role is enabled or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceIncidentRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Incident Role: %s", name))

	s := &client.IncidentRole{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("summary"); ok {
		s.Summary = value.(string)
	}

	if v, ok := d.GetOk("enabled"); ok {
		s.Enabled = tools.Bool(v.(bool))
	}

	res, err := c.CreateIncidentRole(s)
	if err != nil {
		return diag.Errorf("Error creating incident role: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an incident role resource: %v (%s)", name, d.Id()))

	return resourceIncidentRoleRead(ctx, d, meta)
}

func resourceIncidentRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Incident Role: %s", d.Id()))

	res, err := c.GetIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident role: %s", d.Id())
	}

	d.Set("name", res.Name)
	d.Set("description", res.Description)
	d.Set("summary", res.Summary)
	d.Set("enabled", res.Enabled)

	return nil
}

func resourceIncidentRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Incident Role: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.IncidentRole{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("summary") {
		s.Summary = d.Get("summary").(string)
	}

	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateIncidentRole(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident role: %s", err.Error())
	}

	return resourceIncidentRoleRead(ctx, d, meta)
}

func resourceIncidentRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Incident Role: %s", d.Id()))

	err := c.DeleteIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident role: %s", err.Error())
	}

	d.SetId("")

	return nil
}
