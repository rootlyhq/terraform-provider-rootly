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

func resourceHeartbeat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHeartbeatCreate,
		ReadContext:   resourceHeartbeatRead,
		UpdateContext: resourceHeartbeatUpdate,
		DeleteContext: resourceHeartbeatDelete,
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
				Description: "The name of the heartbeat",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the heartbeat",
			},

			"alert_summary": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Summary of alerts triggered when heartbeat expires.",
			},

			"interval": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "",
			},

			"interval_unit": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "seconds",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `seconds`, `minutes`, `hours`.",
			},

			"notification_target_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "",
			},

			"notification_target_type": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "User",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `User`, `Group`, `Service`, `EscalationPolicy`.",
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"status": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "waiting",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `waiting`, `active`, `expired`.",
			},

			"last_pinged_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "When the heartbeat was last pinged.",
			},

			"expires_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "When heartbeat expires",
			},
		},
	}
}

func resourceHeartbeatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Heartbeat"))

	s := &client.Heartbeat{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("alert_summary"); ok {
		s.AlertSummary = value.(string)
	}
	if value, ok := d.GetOkExists("interval"); ok {
		s.Interval = value.(int)
	}
	if value, ok := d.GetOkExists("interval_unit"); ok {
		s.IntervalUnit = value.(string)
	}
	if value, ok := d.GetOkExists("notification_target_id"); ok {
		s.NotificationTargetId = value.(string)
	}
	if value, ok := d.GetOkExists("notification_target_type"); ok {
		s.NotificationTargetType = value.(string)
	}
	if value, ok := d.GetOkExists("enabled"); ok {
		s.Enabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("status"); ok {
		s.Status = value.(string)
	}
	if value, ok := d.GetOkExists("last_pinged_at"); ok {
		s.LastPingedAt = value.(string)
	}
	if value, ok := d.GetOkExists("expires_at"); ok {
		s.ExpiresAt = value.(string)
	}

	res, err := c.CreateHeartbeat(s)
	if err != nil {
		return diag.Errorf("Error creating heartbeat: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a heartbeat resource: %s", d.Id()))

	return resourceHeartbeatRead(ctx, d, meta)
}

func resourceHeartbeatRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Heartbeat: %s", d.Id()))

	item, err := c.GetHeartbeat(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Heartbeat (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading heartbeat: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("alert_summary", item.AlertSummary)
	d.Set("interval", item.Interval)
	d.Set("interval_unit", item.IntervalUnit)
	d.Set("notification_target_id", item.NotificationTargetId)
	d.Set("notification_target_type", item.NotificationTargetType)
	d.Set("enabled", item.Enabled)
	d.Set("status", item.Status)
	d.Set("last_pinged_at", item.LastPingedAt)
	d.Set("expires_at", item.ExpiresAt)

	return nil
}

func resourceHeartbeatUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Heartbeat: %s", d.Id()))

	s := &client.Heartbeat{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("alert_summary") {
		s.AlertSummary = d.Get("alert_summary").(string)
	}
	if d.HasChange("interval") {
		s.Interval = d.Get("interval").(int)
	}
	if d.HasChange("interval_unit") {
		s.IntervalUnit = d.Get("interval_unit").(string)
	}
	if d.HasChange("notification_target_id") {
		s.NotificationTargetId = d.Get("notification_target_id").(string)
	}
	if d.HasChange("notification_target_type") {
		s.NotificationTargetType = d.Get("notification_target_type").(string)
	}
	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}
	if d.HasChange("status") {
		s.Status = d.Get("status").(string)
	}
	if d.HasChange("last_pinged_at") {
		s.LastPingedAt = d.Get("last_pinged_at").(string)
	}
	if d.HasChange("expires_at") {
		s.ExpiresAt = d.Get("expires_at").(string)
	}

	_, err := c.UpdateHeartbeat(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating heartbeat: %s", err.Error())
	}

	return resourceHeartbeatRead(ctx, d, meta)
}

func resourceHeartbeatDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Heartbeat: %s", d.Id()))

	err := c.DeleteHeartbeat(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Heartbeat (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting heartbeat: %s", err.Error())
	}

	d.SetId("")

	return nil
}