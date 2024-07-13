package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceOnCallShadow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOnCallShadowCreate,
		ReadContext:   resourceOnCallShadowRead,
		UpdateContext: resourceOnCallShadowUpdate,
		DeleteContext: resourceOnCallShadowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"schedule_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of schedule the shadow shift belongs to",
			},

			"shadowable_type": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "User",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `User`, `Schedule`.",
			},

			"shadowable_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "ID of schedule or user the shadow user is shadowing",
			},

			"shadow_user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Which user the shadow shift belongs to.",
			},

			"starts_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Start datetime of shadow shift",
			},

			"ends_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "End datetime for shadow shift",
			},
		},
	}
}

func resourceOnCallShadowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating OnCallShadow"))

	s := &client.OnCallShadow{}

	if value, ok := d.GetOkExists("schedule_id"); ok {
		s.ScheduleId = value.(string)
	}
	if value, ok := d.GetOkExists("shadowable_type"); ok {
		s.ShadowableType = value.(string)
	}
	if value, ok := d.GetOkExists("shadowable_id"); ok {
		s.ShadowableId = value.(string)
	}
	if value, ok := d.GetOkExists("shadow_user_id"); ok {
		s.ShadowUserId = value.(int)
	}
	if value, ok := d.GetOkExists("starts_at"); ok {
		s.StartsAt = value.(string)
	}
	if value, ok := d.GetOkExists("ends_at"); ok {
		s.EndsAt = value.(string)
	}

	res, err := c.CreateOnCallShadow(s)
	if err != nil {
		return diag.Errorf("Error creating on_call_shadow: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a on_call_shadow resource: %s", d.Id()))

	return resourceOnCallShadowRead(ctx, d, meta)
}

func resourceOnCallShadowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading OnCallShadow: %s", d.Id()))

	item, err := c.GetOnCallShadow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("OnCallShadow (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading on_call_shadow: %s", d.Id())
	}

	d.Set("schedule_id", item.ScheduleId)
	d.Set("shadowable_type", item.ShadowableType)
	d.Set("shadowable_id", item.ShadowableId)
	d.Set("shadow_user_id", item.ShadowUserId)
	d.Set("starts_at", item.StartsAt)
	d.Set("ends_at", item.EndsAt)

	return nil
}

func resourceOnCallShadowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating OnCallShadow: %s", d.Id()))

	s := &client.OnCallShadow{}

	if d.HasChange("schedule_id") {
		s.ScheduleId = d.Get("schedule_id").(string)
	}
	if d.HasChange("shadowable_type") {
		s.ShadowableType = d.Get("shadowable_type").(string)
	}
	if d.HasChange("shadowable_id") {
		s.ShadowableId = d.Get("shadowable_id").(string)
	}
	if d.HasChange("shadow_user_id") {
		s.ShadowUserId = d.Get("shadow_user_id").(int)
	}
	if d.HasChange("starts_at") {
		s.StartsAt = d.Get("starts_at").(string)
	}
	if d.HasChange("ends_at") {
		s.EndsAt = d.Get("ends_at").(string)
	}

	_, err := c.UpdateOnCallShadow(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating on_call_shadow: %s", err.Error())
	}

	return resourceOnCallShadowRead(ctx, d, meta)
}

func resourceOnCallShadowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting OnCallShadow: %s", d.Id()))

	err := c.DeleteOnCallShadow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("OnCallShadow (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting on_call_shadow: %s", err.Error())
	}

	d.SetId("")

	return nil
}
