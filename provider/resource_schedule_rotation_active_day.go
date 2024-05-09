package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceScheduleRotationActiveDay() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleRotationActiveDayCreate,
		ReadContext:   resourceScheduleRotationActiveDayRead,
		UpdateContext: resourceScheduleRotationActiveDayUpdate,
		DeleteContext: resourceScheduleRotationActiveDayDelete,
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

			"day_name": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "S",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Schedule rotation day name for which active times to be created. Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.",
			},

			"active_time_attributes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Schedule rotation active times per day",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceScheduleRotationActiveDayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ScheduleRotationActiveDay"))

	s := &client.ScheduleRotationActiveDay{}

	if value, ok := d.GetOkExists("schedule_rotation_id"); ok {
		s.ScheduleRotationId = value.(string)
	}
	if value, ok := d.GetOkExists("day_name"); ok {
		s.DayName = value.(string)
	}
	if value, ok := d.GetOkExists("active_time_attributes"); ok {
		s.ActiveTimeAttributes = value.([]interface{})
	}

	res, err := c.CreateScheduleRotationActiveDay(s)
	if err != nil {
		return diag.Errorf("Error creating schedule_rotation_active_day: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a schedule_rotation_active_day resource: %s", d.Id()))

	return resourceScheduleRotationActiveDayRead(ctx, d, meta)
}

func resourceScheduleRotationActiveDayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ScheduleRotationActiveDay: %s", d.Id()))

	item, err := c.GetScheduleRotationActiveDay(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotationActiveDay (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading schedule_rotation_active_day: %s", d.Id())
	}

	d.Set("schedule_rotation_id", item.ScheduleRotationId)
	d.Set("day_name", item.DayName)
	d.Set("active_time_attributes", item.ActiveTimeAttributes)

	return nil
}

func resourceScheduleRotationActiveDayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ScheduleRotationActiveDay: %s", d.Id()))

	s := &client.ScheduleRotationActiveDay{}

	if d.HasChange("schedule_rotation_id") {
		s.ScheduleRotationId = d.Get("schedule_rotation_id").(string)
	}
	if d.HasChange("day_name") {
		s.DayName = d.Get("day_name").(string)
	}
	if d.HasChange("active_time_attributes") {
		s.ActiveTimeAttributes = d.Get("active_time_attributes").([]interface{})
	}

	_, err := c.UpdateScheduleRotationActiveDay(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating schedule_rotation_active_day: %s", err.Error())
	}

	return resourceScheduleRotationActiveDayRead(ctx, d, meta)
}

func resourceScheduleRotationActiveDayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ScheduleRotationActiveDay: %s", d.Id()))

	err := c.DeleteScheduleRotationActiveDay(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotationActiveDay (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting schedule_rotation_active_day: %s", err.Error())
	}

	d.SetId("")

	return nil
}
