package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceAlertUrgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertUrgencyCreate,
		ReadContext:   resourceAlertUrgencyRead,
		UpdateContext: resourceAlertUrgencyUpdate,
		DeleteContext: resourceAlertUrgencyDelete,
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
				Description: "The name of the alert urgency",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The description of the alert urgency",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Position of the alert urgency",
			},
		},
	}
}

func resourceAlertUrgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating AlertUrgency"))

	s := &client.AlertUrgency{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateAlertUrgency(s)
	if err != nil {
		return diag.Errorf("Error creating alert_urgency: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a alert_urgency resource: %s", d.Id()))

	return resourceAlertUrgencyRead(ctx, d, meta)
}

func resourceAlertUrgencyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading AlertUrgency: %s", d.Id()))

	item, err := c.GetAlertUrgency(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertUrgency (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading alert_urgency: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("position", item.Position)

	return nil
}

func resourceAlertUrgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating AlertUrgency: %s", d.Id()))

	s := &client.AlertUrgency{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateAlertUrgency(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating alert_urgency: %s", err.Error())
	}

	return resourceAlertUrgencyRead(ctx, d, meta)
}

func resourceAlertUrgencyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting AlertUrgency: %s", d.Id()))

	err := c.DeleteAlertUrgency(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertUrgency (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting alert_urgency: %s", err.Error())
	}

	d.SetId("")

	return nil
}