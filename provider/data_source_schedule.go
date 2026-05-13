package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/polling"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Filter by date range using 'lt' and 'gt'.",
				Optional:    true,
			},

			"sync_linear_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"include_shadows_in_slack_notifications": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"shift_start_notifications_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"shift_update_notifications_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"shift_report_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"shift_report_day_of_week": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"shift_report_time_of_day": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"shift_report_time_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListSchedulesParams)
	page_size := 1
	params.PageSize = &page_size

	if value, ok := d.GetOkExists("name"); ok {
		name := value.(string)
		params.FilterName = &name
	}

	created_at_gt := d.Get("created_at").(map[string]interface{})
	if value, exists := created_at_gt["gt"]; exists {
		v := value.(string)
		params.FilterCreatedAtGt = &v
	}

	created_at_lt := d.Get("created_at").(map[string]interface{})
	if value, exists := created_at_lt["lt"]; exists {
		v := value.(string)
		params.FilterCreatedAtLt = &v
	}

	items, err := polling.WaitForList(ctx, "schedule", func() ([]interface{}, error) {
		return c.ListSchedules(params)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	item, _ := items[0].(*client.Schedule)

	d.SetId(item.ID)
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
