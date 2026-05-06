package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/provider/stateupgrade"
	"github.com/rootlyhq/terraform-provider-rootly/v5/tools"
)

func resourceSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    stateupgrade.ScheduleV0().CoreConfigSchema().ImpliedType(),
				Upgrade: stateupgrade.UpgradeScheduleV0ToV1,
			},
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the schedule",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the schedule",
			},

			"all_time_coverage": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "24/7 coverage of the schedule. Value must be one of true or false",
			},

			"slack_user_group": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default:     map[string]interface{}{},
				Optional:    true,
				Description: "Map must contain two fields, `id` and `name`. Synced slack group of the schedule",
			},

			"slack_channel": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default:     map[string]interface{}{},
				Optional:    true,
				Description: "Map must contain two fields, `id` and `name`. Synced slack channel of the schedule",
			},

			"owner_group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Required:         false,
				Optional:         true,
				Description:      "The owning teams for this schedules.",
			},

			"owner_user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "ID of user assigned as owner of the schedule. Defaults to the API token's user if not specified.",
			},

			"sync_linear_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether the schedule is synced with Linear. Value must be one of true or false",
			},

			"include_shadows_in_slack_notifications": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether shadow users are included in Slack notifications and user group syncing. Value must be one of true or false",
			},

			"shift_start_notifications_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether to send a Slack message every time a new shift begins. Requires `slack_channel` to be set. Value must be one of true or false",
			},

			"shift_update_notifications_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether to send a Slack message whenever a shift is updated (overrides, removed users, rotation changes, etc.). Requires `slack_channel` to be set. Value must be one of true or false",
			},

			"shift_report_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether the weekly shift summary report is sent. Requires `slack_channel` to be set. Value must be one of true or false",
			},

			"shift_report_day_of_week": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Day of week the weekly shift summary is sent. Value must be one of `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`.",
			},

			"shift_report_time_of_day": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Time of day the weekly shift summary is sent, in `HH:MM` 24-hour format.",
			},

			"shift_report_time_zone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "IANA time zone used for the weekly shift summary (e.g. `Australia/Sydney`).",
			},
		},
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Schedule"))

	s := &client.Schedule{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("all_time_coverage"); ok {
		s.AllTimeCoverage = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("slack_user_group"); ok {
		slackUserGroup := value.(map[string]interface{})
		if len(slackUserGroup) == 0 {
			s.SlackUserGroup = map[string]interface{}{}
		} else {
			s.SlackUserGroup = slackUserGroup
		}
	} else {
		s.SlackUserGroup = map[string]interface{}{}
	}
	if value, ok := d.GetOkExists("slack_channel"); ok {
		slackChannel := value.(map[string]interface{})
		if len(slackChannel) == 0 {
			s.SlackChannel = map[string]interface{}{}
		} else {
			s.SlackChannel = slackChannel
		}
	} else {
		s.SlackChannel = map[string]interface{}{}
	}
	if value, ok := d.GetOkExists("owner_group_ids"); ok {
		s.OwnerGroupIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("owner_user_id"); ok {
		s.OwnerUserId = value.(int)
	}
	if value, ok := d.GetOkExists("sync_linear_enabled"); ok {
		s.SyncLinearEnabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("include_shadows_in_slack_notifications"); ok {
		s.IncludeShadowsInSlackNotifications = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("shift_start_notifications_enabled"); ok {
		s.ShiftStartNotificationsEnabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("shift_update_notifications_enabled"); ok {
		s.ShiftUpdateNotificationsEnabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("shift_report_enabled"); ok {
		s.ShiftReportEnabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("shift_report_day_of_week"); ok {
		s.ShiftReportDayOfWeek = value.(string)
	}
	if value, ok := d.GetOkExists("shift_report_time_of_day"); ok {
		s.ShiftReportTimeOfDay = value.(string)
	}
	if value, ok := d.GetOkExists("shift_report_time_zone"); ok {
		s.ShiftReportTimeZone = value.(string)
	}

	res, err := c.CreateSchedule(s)
	if err != nil {
		return diag.Errorf("Error creating schedule: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a schedule resource: %s", d.Id()))

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Schedule: %s", d.Id()))

	item, err := c.GetSchedule(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Schedule (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading schedule: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("all_time_coverage", item.AllTimeCoverage)
	d.Set("slack_user_group", item.SlackUserGroup)
	d.Set("slack_channel", item.SlackChannel)
	d.Set("owner_group_ids", item.OwnerGroupIds)
	d.Set("owner_user_id", item.OwnerUserId)
	d.Set("sync_linear_enabled", item.SyncLinearEnabled)
	d.Set("include_shadows_in_slack_notifications", item.IncludeShadowsInSlackNotifications)
	d.Set("shift_start_notifications_enabled", item.ShiftStartNotificationsEnabled)
	d.Set("shift_update_notifications_enabled", item.ShiftUpdateNotificationsEnabled)
	d.Set("shift_report_enabled", item.ShiftReportEnabled)
	d.Set("shift_report_day_of_week", item.ShiftReportDayOfWeek)
	d.Set("shift_report_time_of_day", item.ShiftReportTimeOfDay)
	d.Set("shift_report_time_zone", item.ShiftReportTimeZone)

	return nil
}

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Schedule: %s", d.Id()))

	s := &client.Schedule{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("all_time_coverage") {
		s.AllTimeCoverage = tools.Bool(d.Get("all_time_coverage").(bool))
	}
	if d.HasChange("slack_user_group") {
		if value := d.Get("slack_user_group"); value != nil {
			slackUserGroup := value.(map[string]interface{})
			if len(slackUserGroup) == 0 {
				s.SlackUserGroup = map[string]interface{}{}
			} else {
				s.SlackUserGroup = slackUserGroup
			}
		} else {
			s.SlackUserGroup = map[string]interface{}{}
		}
	}
	if d.HasChange("slack_channel") {
		if value := d.Get("slack_channel"); value != nil {
			slackChannel := value.(map[string]interface{})
			if len(slackChannel) == 0 {
				s.SlackChannel = map[string]interface{}{}
			} else {
				s.SlackChannel = slackChannel
			}
		} else {
			s.SlackChannel = map[string]interface{}{}
		}
	}
	if d.HasChange("owner_group_ids") {
		s.OwnerGroupIds = d.Get("owner_group_ids").([]interface{})
	}
	if d.HasChange("owner_user_id") {
		s.OwnerUserId = d.Get("owner_user_id").(int)
	}
	if d.HasChange("sync_linear_enabled") {
		s.SyncLinearEnabled = tools.Bool(d.Get("sync_linear_enabled").(bool))
	}
	if d.HasChange("include_shadows_in_slack_notifications") {
		s.IncludeShadowsInSlackNotifications = tools.Bool(d.Get("include_shadows_in_slack_notifications").(bool))
	}
	if d.HasChange("shift_start_notifications_enabled") {
		s.ShiftStartNotificationsEnabled = tools.Bool(d.Get("shift_start_notifications_enabled").(bool))
	}
	if d.HasChange("shift_update_notifications_enabled") {
		s.ShiftUpdateNotificationsEnabled = tools.Bool(d.Get("shift_update_notifications_enabled").(bool))
	}
	if d.HasChange("shift_report_enabled") {
		s.ShiftReportEnabled = tools.Bool(d.Get("shift_report_enabled").(bool))
	}
	if d.HasChange("shift_report_day_of_week") {
		s.ShiftReportDayOfWeek = d.Get("shift_report_day_of_week").(string)
	}
	if d.HasChange("shift_report_time_of_day") {
		s.ShiftReportTimeOfDay = d.Get("shift_report_time_of_day").(string)
	}
	if d.HasChange("shift_report_time_zone") {
		s.ShiftReportTimeZone = d.Get("shift_report_time_zone").(string)
	}

	_, err := c.UpdateSchedule(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating schedule: %s", err.Error())
	}

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Schedule: %s", d.Id()))

	err := c.DeleteSchedule(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Schedule (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting schedule: %s", err.Error())
	}

	d.SetId("")

	return nil
}
