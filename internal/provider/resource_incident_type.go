package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceIncidentType() *schema.Resource {
	return &schema.Resource{
		Description: "Manages incident types (e.g Cloud, Customer Facing, Security, Training).",

		CreateContext: resourceIncidentTypeCreate,
		ReadContext:   resourceIncidentTypeRead,
		UpdateContext: resourceIncidentTypeUpdate,
		DeleteContext: resourceIncidentTypeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the incident type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the incident type",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The cikir of the incident type",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#047BF8", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
		},
	}
}

func resourceIncidentTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Incident Type: %s", name))

	s := &client.IncidentType{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if value, ok := d.GetOk("color"); ok {
		s.Color = value.(string)
	}

	res, err := c.CreateIncidentType(s)
	if err != nil {
		return diag.Errorf("Error creating incident type: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an incident type resource: %v (%s)", name, d.Id()))

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Incident Type: %s", d.Id()))

	res, err := c.GetIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident type: %s", d.Id())
	}

	d.Set("name", res.Name)
	d.Set("description", res.Description)
	d.Set("color", res.Color)

	return nil
}

func resourceIncidentTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Incident Type: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.IncidentType{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}

	_, err := c.UpdateIncidentType(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident type: %s", err.Error())
	}

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Incident Type: %s", d.Id()))

	err := c.DeleteIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident type: %s", err.Error())
	}

	d.SetId("")

	return nil
}
