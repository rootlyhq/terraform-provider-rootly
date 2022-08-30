package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceIncidentType() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceIncidentTypeCreate,
		ReadContext: resourceIncidentTypeRead,
		UpdateContext: resourceIncidentTypeUpdate,
		DeleteContext: resourceIncidentTypeDelete,
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
				Description: "The name of the incident type",
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the incident type",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the incident type",
			},
			

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			
		},
	}
}

func resourceIncidentTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentType"))

	s := &client.IncidentType{}

	  if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
    if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
    if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
    if value, ok := d.GetOkExists("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateIncidentType(s)
	if err != nil {
		return diag.Errorf("Error creating incident_type: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_type resource: %s", d.Id()))

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentType: %s", d.Id()))

	item, err := c.GetIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_type: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("color", item.Color)

	return nil
}

func resourceIncidentTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentType: %s", d.Id()))

	s := &client.IncidentType{}

	  if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
    if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
    if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
    if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	_, err := c.UpdateIncidentType(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_type: %s", err.Error())
	}

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentType: %s", d.Id()))

	err := c.DeleteIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_type: %s", err.Error())
	}

	d.SetId("")

	return nil
}
