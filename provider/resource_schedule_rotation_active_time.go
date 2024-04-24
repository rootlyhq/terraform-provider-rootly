package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceScheduleRotationActiveTime() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleRotationActiveTimeCreate,
		ReadContext:   resourceScheduleRotationActiveTimeRead,
		UpdateContext: resourceScheduleRotationActiveTimeUpdate,
		DeleteContext: resourceScheduleRotationActiveTimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"schedule_rotation_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},

			"start_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Start time for schedule rotation active time",
			},

			"end_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "End time for schedule rotation active time",
			},
		},
	}
}

func resourceScheduleRotationActiveTimeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ScheduleRotationActiveTime"))

	s := &client.ScheduleRotationActiveTime{}

	if value, ok := d.GetOkExists("schedule_rotation_id"); ok {
		s.ScheduleRotationId = value.(string)
	}
	if value, ok := d.GetOkExists("start_time"); ok {
		s.StartTime = value.(string)
	}
	if value, ok := d.GetOkExists("end_time"); ok {
		s.EndTime = value.(string)
	}

	res, err := c.CreateScheduleRotationActiveTime(s)
	if err != nil {
		return diag.Errorf("Error creating schedule_rotation_active_time: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a schedule_rotation_active_time resource: %s", d.Id()))

	return resourceScheduleRotationActiveTimeRead(ctx, d, meta)
}

func resourceScheduleRotationActiveTimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ScheduleRotationActiveTime: %s", d.Id()))

	item, err := c.GetScheduleRotationActiveTime(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotationActiveTime (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading schedule_rotation_active_time: %s", d.Id())
	}

	d.Set("schedule_rotation_id", item.ScheduleRotationId)
	d.Set("start_time", item.StartTime)
	d.Set("end_time", item.EndTime)

	return nil
}

func resourceScheduleRotationActiveTimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ScheduleRotationActiveTime: %s", d.Id()))

	s := &client.ScheduleRotationActiveTime{}

	if d.HasChange("schedule_rotation_id") {
		s.ScheduleRotationId = d.Get("schedule_rotation_id").(string)
	}
	if d.HasChange("start_time") {
		s.StartTime = d.Get("start_time").(string)
	}
	if d.HasChange("end_time") {
		s.EndTime = d.Get("end_time").(string)
	}

	_, err := c.UpdateScheduleRotationActiveTime(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating schedule_rotation_active_time: %s", err.Error())
	}

	return resourceScheduleRotationActiveTimeRead(ctx, d, meta)
}

func resourceScheduleRotationActiveTimeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ScheduleRotationActiveTime: %s", d.Id()))

	err := c.DeleteScheduleRotationActiveTime(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotationActiveTime (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting schedule_rotation_active_time: %s", err.Error())
	}

	d.SetId("")

	return nil
}
