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

func resourceOverrideShift() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOverrideShiftCreate,
		ReadContext:   resourceOverrideShiftRead,
		UpdateContext: resourceOverrideShiftUpdate,
		DeleteContext: resourceOverrideShiftDelete,
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
				Description: "ID of schedule",
			},

			"rotation_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "ID of rotation",
			},

			"user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "User to whom override shift is assigned to",
			},

			"starts_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Start datetime of shift",
			},

			"ends_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "End datetime of shift",
			},

			"is_override": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Denotes shift is an override shift. Value must be one of true or false",
			},

			"shift_override": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Override metadata",
			},
		},
	}
}

func resourceOverrideShiftCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating OverrideShift"))

	s := &client.OverrideShift{}

	if value, ok := d.GetOkExists("schedule_id"); ok {
		s.ScheduleId = value.(string)
	}
	if value, ok := d.GetOkExists("rotation_id"); ok {
		s.RotationId = value.(string)
	}
	if value, ok := d.GetOkExists("user_id"); ok {
		s.UserId = value.(int)
	}
	if value, ok := d.GetOkExists("starts_at"); ok {
		s.StartsAt = value.(string)
	}
	if value, ok := d.GetOkExists("ends_at"); ok {
		s.EndsAt = value.(string)
	}
	if value, ok := d.GetOkExists("is_override"); ok {
		s.IsOverride = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("shift_override"); ok {
		s.ShiftOverride = value.(map[string]interface{})
	}

	res, err := c.CreateOverrideShift(s)
	if err != nil {
		return diag.Errorf("Error creating override_shift: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a override_shift resource: %s", d.Id()))

	return resourceOverrideShiftRead(ctx, d, meta)
}

func resourceOverrideShiftRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading OverrideShift: %s", d.Id()))

	item, err := c.GetOverrideShift(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("OverrideShift (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading override_shift: %s", d.Id())
	}

	d.Set("schedule_id", item.ScheduleId)
	d.Set("rotation_id", item.RotationId)
	d.Set("user_id", item.UserId)
	d.Set("starts_at", item.StartsAt)
	d.Set("ends_at", item.EndsAt)
	d.Set("is_override", item.IsOverride)
	d.Set("shift_override", item.ShiftOverride)

	return nil
}

func resourceOverrideShiftUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating OverrideShift: %s", d.Id()))

	s := &client.OverrideShift{}

	if d.HasChange("schedule_id") {
		s.ScheduleId = d.Get("schedule_id").(string)
	}
	if d.HasChange("rotation_id") {
		s.RotationId = d.Get("rotation_id").(string)
	}
	if d.HasChange("user_id") {
		s.UserId = d.Get("user_id").(int)
	}
	if d.HasChange("starts_at") {
		s.StartsAt = d.Get("starts_at").(string)
	}
	if d.HasChange("ends_at") {
		s.EndsAt = d.Get("ends_at").(string)
	}
	if d.HasChange("is_override") {
		s.IsOverride = tools.Bool(d.Get("is_override").(bool))
	}
	if d.HasChange("shift_override") {
		s.ShiftOverride = d.Get("shift_override").(map[string]interface{})
	}

	_, err := c.UpdateOverrideShift(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating override_shift: %s", err.Error())
	}

	return resourceOverrideShiftRead(ctx, d, meta)
}

func resourceOverrideShiftDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting OverrideShift: %s", d.Id()))

	err := c.DeleteOverrideShift(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("OverrideShift (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting override_shift: %s", err.Error())
	}

	d.SetId("")

	return nil
}
