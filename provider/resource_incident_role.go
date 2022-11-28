package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceIncidentRole() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceIncidentRoleCreate,
		ReadContext: resourceIncidentRoleRead,
		UpdateContext: resourceIncidentRoleUpdate,
		DeleteContext: resourceIncidentRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the incident role",
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the incident role",
			},
			

			"summary": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The summary of the incident role",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the incident role",
			},
			
		},
	}
}


func resourceIncidentRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentRole"))

	s := &client.IncidentRole{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("summary"); ok {
				s.Summary = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}

	res, err := c.CreateIncidentRole(s)
	if err != nil {
		return diag.Errorf("Error creating incident_role: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_role resource: %s", d.Id()))

	return resourceIncidentRoleRead(ctx, d, meta)
}


func resourceIncidentRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentRole: %s", d.Id()))

	item, err := c.GetIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_role: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("summary", item.Summary)
  d.Set("description", item.Description)

	return nil
}


func resourceIncidentRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentRole: %s", d.Id()))

	s := &client.IncidentRole{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("summary") {
				s.Summary = d.Get("summary").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}

	_, err := c.UpdateIncidentRole(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_role: %s", err.Error())
	}

	return resourceIncidentRoleRead(ctx, d, meta)
}


func resourceIncidentRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentRole: %s", d.Id()))

	err := c.DeleteIncidentRole(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRole (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_role: %s", err.Error())
	}

	d.SetId("")

	return nil
}

