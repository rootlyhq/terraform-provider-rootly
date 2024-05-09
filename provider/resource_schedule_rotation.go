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

func resourceScheduleRotation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleRotationCreate,
		ReadContext:   resourceScheduleRotationRead,
		UpdateContext: resourceScheduleRotationUpdate,
		DeleteContext: resourceScheduleRotationDelete,
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
				Description: "The ID of parent schedule",
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the schedule rotation",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Position of the schedule rotation",
			},

			"schedule_rotationable_type": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "ScheduleDailyRotation",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Schedule rotation type. Value must be one of `ScheduleDailyRotation`, `ScheduleWeeklyRotation`, `ScheduleBiweeklyRotation`, `ScheduleMonthlyRotation`, `ScheduleCustomRotation`.",
			},

			"active_all_week": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Schedule rotation active all week?. Value must be one of true or false",
			},

			"active_days": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.",
			},

			"active_time_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"active_time_attributes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Schedule rotation's active times",
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

			"time_zone": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "UTC",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Schedule Rotation Timezone. Value must be one of `International Date Line West`, `American Samoa`, `Midway Island`, `Hawaii`, `Alaska`, `Pacific Time (US & Canada)`, `Tijuana`, `Arizona`, `Mazatlan`, `Mountain Time (US & Canada)`, `Central America`, `Central Time (US & Canada)`, `Chihuahua`, `Guadalajara`, `Mexico City`, `Monterrey`, `Saskatchewan`, `Bogota`, `Eastern Time (US & Canada)`, `Indiana (East)`, `Lima`, `Quito`, `Atlantic Time (Canada)`, `Caracas`, `Georgetown`, `La Paz`, `Puerto Rico`, `Santiago`, `Newfoundland`, `Brasilia`, `Buenos Aires`, `Montevideo`, `Greenland`, `Mid-Atlantic`, `Azores`, `Cape Verde Is.`, `Edinburgh`, `Lisbon`, `London`, `Monrovia`, `UTC`, `Amsterdam`, `Belgrade`, `Berlin`, `Bern`, `Bratislava`, `Brussels`, `Budapest`, `Casablanca`, `Copenhagen`, `Dublin`, `Ljubljana`, `Madrid`, `Paris`, `Prague`, `Rome`, `Sarajevo`, `Skopje`, `Stockholm`, `Vienna`, `Warsaw`, `West Central Africa`, `Zagreb`, `Zurich`, `Athens`, `Bucharest`, `Cairo`, `Harare`, `Helsinki`, `Jerusalem`, `Kaliningrad`, `Kyiv`, `Pretoria`, `Riga`, `Sofia`, `Tallinn`, `Vilnius`, `Baghdad`, `Istanbul`, `Kuwait`, `Minsk`, `Moscow`, `Nairobi`, `Riyadh`, `St. Petersburg`, `Volgograd`, `Tehran`, `Abu Dhabi`, `Baku`, `Muscat`, `Samara`, `Tbilisi`, `Yerevan`, `Kabul`, `Almaty`, `Ekaterinburg`, `Islamabad`, `Karachi`, `Tashkent`, `Chennai`, `Kolkata`, `Mumbai`, `New Delhi`, `Sri Jayawardenepura`, `Kathmandu`, `Astana`, `Dhaka`, `Urumqi`, `Rangoon`, `Bangkok`, `Hanoi`, `Jakarta`, `Krasnoyarsk`, `Novosibirsk`, `Beijing`, `Chongqing`, `Hong Kong`, `Irkutsk`, `Kuala Lumpur`, `Perth`, `Singapore`, `Taipei`, `Ulaanbaatar`, `Osaka`, `Sapporo`, `Seoul`, `Tokyo`, `Yakutsk`, `Adelaide`, `Darwin`, `Brisbane`, `Canberra`, `Guam`, `Hobart`, `Melbourne`, `Port Moresby`, `Sydney`, `Vladivostok`, `Magadan`, `New Caledonia`, `Solomon Is.`, `Srednekolymsk`, `Auckland`, `Fiji`, `Kamchatka`, `Marshall Is.`, `Wellington`, `Chatham Is.`, `Nuku'alofa`, `Samoa`, `Tokelau Is.`.",
			},

			"schedule_rotationable_attributes": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "",
			},
		},
	}
}

func resourceScheduleRotationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ScheduleRotation"))

	s := &client.ScheduleRotation{}

	if value, ok := d.GetOkExists("schedule_id"); ok {
		s.ScheduleId = value.(string)
	}
	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}
	if value, ok := d.GetOkExists("schedule_rotationable_type"); ok {
		s.ScheduleRotationableType = value.(string)
	}
	if value, ok := d.GetOkExists("active_all_week"); ok {
		s.ActiveAllWeek = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("active_days"); ok {
		s.ActiveDays = value.([]interface{})
	}
	if value, ok := d.GetOkExists("active_time_type"); ok {
		s.ActiveTimeType = value.(string)
	}
	if value, ok := d.GetOkExists("active_time_attributes"); ok {
		s.ActiveTimeAttributes = value.([]interface{})
	}
	if value, ok := d.GetOkExists("time_zone"); ok {
		s.TimeZone = value.(string)
	}
	if value, ok := d.GetOkExists("schedule_rotationable_attributes"); ok {
		s.ScheduleRotationableAttributes = value.(map[string]interface{})
	}

	res, err := c.CreateScheduleRotation(s)
	if err != nil {
		return diag.Errorf("Error creating schedule_rotation: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a schedule_rotation resource: %s", d.Id()))

	return resourceScheduleRotationRead(ctx, d, meta)
}

func resourceScheduleRotationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ScheduleRotation: %s", d.Id()))

	item, err := c.GetScheduleRotation(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotation (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading schedule_rotation: %s", d.Id())
	}

	d.Set("schedule_id", item.ScheduleId)
	d.Set("name", item.Name)
	d.Set("position", item.Position)
	d.Set("schedule_rotationable_type", item.ScheduleRotationableType)
	d.Set("active_all_week", item.ActiveAllWeek)
	d.Set("active_days", item.ActiveDays)
	d.Set("active_time_type", item.ActiveTimeType)
	d.Set("active_time_attributes", item.ActiveTimeAttributes)
	d.Set("time_zone", item.TimeZone)
	d.Set("schedule_rotationable_attributes", item.ScheduleRotationableAttributes)

	return nil
}

func resourceScheduleRotationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ScheduleRotation: %s", d.Id()))

	s := &client.ScheduleRotation{}

	if d.HasChange("schedule_id") {
		s.ScheduleId = d.Get("schedule_id").(string)
	}
	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}
	if d.HasChange("schedule_rotationable_type") {
		s.ScheduleRotationableType = d.Get("schedule_rotationable_type").(string)
	}
	if d.HasChange("active_all_week") {
		s.ActiveAllWeek = tools.Bool(d.Get("active_all_week").(bool))
	}
	if d.HasChange("active_days") {
		s.ActiveDays = d.Get("active_days").([]interface{})
	}
	if d.HasChange("active_time_type") {
		s.ActiveTimeType = d.Get("active_time_type").(string)
	}
	if d.HasChange("active_time_attributes") {
		s.ActiveTimeAttributes = d.Get("active_time_attributes").([]interface{})
	}
	if d.HasChange("time_zone") {
		s.TimeZone = d.Get("time_zone").(string)
	}
	if d.HasChange("schedule_rotationable_attributes") {
		s.ScheduleRotationableAttributes = d.Get("schedule_rotationable_attributes").(map[string]interface{})
	}

	_, err := c.UpdateScheduleRotation(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating schedule_rotation: %s", err.Error())
	}

	return resourceScheduleRotationRead(ctx, d, meta)
}

func resourceScheduleRotationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ScheduleRotation: %s", d.Id()))

	err := c.DeleteScheduleRotation(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("ScheduleRotation (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting schedule_rotation: %s", err.Error())
	}

	d.SetId("")

	return nil
}
