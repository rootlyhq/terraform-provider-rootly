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

func resourceAlertGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertGroupCreate,
		ReadContext:   resourceAlertGroupRead,
		UpdateContext: resourceAlertGroupUpdate,
		DeleteContext: resourceAlertGroupDelete,
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
				Description: "The name of the alert group",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the alert group",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The slug of the alert group",
			},

			"condition_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Grouping condition for the alert group",
			},

			"time_window": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Time window for the alert grouping",
			},

			"group_by_alert_title": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether the alerts are grouped by title or not. Value must be one of true or false",
			},

			"group_by_alert_urgency": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether the alerts are grouped by urgency or not. Value must be one of true or false",
			},

			"deleted_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Date or deletion",
			},
		},
	}
}

func resourceAlertGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating AlertGroup"))

	s := &client.AlertGroup{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("condition_type"); ok {
		s.ConditionType = value.(string)
	}
	if value, ok := d.GetOkExists("time_window"); ok {
		s.TimeWindow = value.(int)
	}
	if value, ok := d.GetOkExists("group_by_alert_title"); ok {
		s.GroupByAlertTitle = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("group_by_alert_urgency"); ok {
		s.GroupByAlertUrgency = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("deleted_at"); ok {
		s.DeletedAt = value.(string)
	}

	res, err := c.CreateAlertGroup(s)
	if err != nil {
		return diag.Errorf("Error creating alert_group: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a alert_group resource: %s", d.Id()))

	return resourceAlertGroupRead(ctx, d, meta)
}

func resourceAlertGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading AlertGroup: %s", d.Id()))

	item, err := c.GetAlertGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading alert_group: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("slug", item.Slug)
	d.Set("condition_type", item.ConditionType)
	d.Set("time_window", item.TimeWindow)
	d.Set("group_by_alert_title", item.GroupByAlertTitle)
	d.Set("group_by_alert_urgency", item.GroupByAlertUrgency)
	d.Set("deleted_at", item.DeletedAt)

	return nil
}

func resourceAlertGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating AlertGroup: %s", d.Id()))

	s := &client.AlertGroup{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("condition_type") {
		s.ConditionType = d.Get("condition_type").(string)
	}
	if d.HasChange("time_window") {
		s.TimeWindow = d.Get("time_window").(int)
	}
	if d.HasChange("group_by_alert_title") {
		s.GroupByAlertTitle = tools.Bool(d.Get("group_by_alert_title").(bool))
	}
	if d.HasChange("group_by_alert_urgency") {
		s.GroupByAlertUrgency = tools.Bool(d.Get("group_by_alert_urgency").(bool))
	}
	if d.HasChange("deleted_at") {
		s.DeletedAt = d.Get("deleted_at").(string)
	}

	_, err := c.UpdateAlertGroup(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating alert_group: %s", err.Error())
	}

	return resourceAlertGroupRead(ctx, d, meta)
}

func resourceAlertGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting AlertGroup: %s", d.Id()))

	err := c.DeleteAlertGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting alert_group: %s", err.Error())
	}

	d.SetId("")

	return nil
}