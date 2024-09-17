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

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
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
				Description: "The name of the team",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the team",
			},

			"notify_emails": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Emails to attach to the team",
			},

			"color": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The hex color of the team",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Position of the team",
			},

			"pagerduty_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The PagerDuty group id associated to this team",
			},

			"pagerduty_service_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The PagerDuty service id associated to this team",
			},

			"opsgenie_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The Opsgenie group id associated to this team",
			},

			"victor_ops_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The VictorOps group id associated to this team",
			},

			"pagertree_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The PagerTree group id associated to this team",
			},

			"cortex_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The Cortex group id associated to this team",
			},

			"service_now_ci_sys_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The Service Now CI sys id associated to this team",
			},

			"user_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "The User ID's members of this team",
			},

			"slack_channels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Slack Channels associated with this service",
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

			"slack_aliases": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Slack Aliases associated with this service",
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

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Team"))

	s := &client.Team{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("notify_emails"); ok {
		s.NotifyEmails = value.([]interface{})
	}
	if value, ok := d.GetOkExists("color"); ok {
		s.Color = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}
	if value, ok := d.GetOkExists("pagerduty_id"); ok {
		s.PagerdutyId = value.(string)
	}
	if value, ok := d.GetOkExists("pagerduty_service_id"); ok {
		s.PagerdutyServiceId = value.(string)
	}
	if value, ok := d.GetOkExists("opsgenie_id"); ok {
		s.OpsgenieId = value.(string)
	}
	if value, ok := d.GetOkExists("victor_ops_id"); ok {
		s.VictorOpsId = value.(string)
	}
	if value, ok := d.GetOkExists("pagertree_id"); ok {
		s.PagertreeId = value.(string)
	}
	if value, ok := d.GetOkExists("cortex_id"); ok {
		s.CortexId = value.(string)
	}
	if value, ok := d.GetOkExists("service_now_ci_sys_id"); ok {
		s.ServiceNowCiSysId = value.(string)
	}
	if value, ok := d.GetOkExists("user_ids"); ok {
		s.UserIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("slack_channels"); ok {
		s.SlackChannels = value.([]interface{})
	}
	if value, ok := d.GetOkExists("slack_aliases"); ok {
		s.SlackAliases = value.([]interface{})
	}

	res, err := c.CreateTeam(s)
	if err != nil {
		return diag.Errorf("Error creating team: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a team resource: %s", d.Id()))

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Team: %s", d.Id()))

	item, err := c.GetTeam(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading team: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("notify_emails", item.NotifyEmails)
	d.Set("color", item.Color)
	d.Set("position", item.Position)
	d.Set("pagerduty_id", item.PagerdutyId)
	d.Set("pagerduty_service_id", item.PagerdutyServiceId)
	d.Set("opsgenie_id", item.OpsgenieId)
	d.Set("victor_ops_id", item.VictorOpsId)
	d.Set("pagertree_id", item.PagertreeId)
	d.Set("cortex_id", item.CortexId)
	d.Set("service_now_ci_sys_id", item.ServiceNowCiSysId)
	d.Set("user_ids", item.UserIds)
	d.Set("slack_channels", item.SlackChannels)
	d.Set("slack_aliases", item.SlackAliases)

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Team: %s", d.Id()))

	s := &client.Team{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("notify_emails") {
		s.NotifyEmails = d.Get("notify_emails").([]interface{})
	}
	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}
	if d.HasChange("pagerduty_id") {
		s.PagerdutyId = d.Get("pagerduty_id").(string)
	}
	if d.HasChange("pagerduty_service_id") {
		s.PagerdutyServiceId = d.Get("pagerduty_service_id").(string)
	}
	if d.HasChange("opsgenie_id") {
		s.OpsgenieId = d.Get("opsgenie_id").(string)
	}
	if d.HasChange("victor_ops_id") {
		s.VictorOpsId = d.Get("victor_ops_id").(string)
	}
	if d.HasChange("pagertree_id") {
		s.PagertreeId = d.Get("pagertree_id").(string)
	}
	if d.HasChange("cortex_id") {
		s.CortexId = d.Get("cortex_id").(string)
	}
	if d.HasChange("service_now_ci_sys_id") {
		s.ServiceNowCiSysId = d.Get("service_now_ci_sys_id").(string)
	}
	if d.HasChange("user_ids") {
		s.UserIds = d.Get("user_ids").([]interface{})
	}
	if d.HasChange("slack_channels") {
		s.SlackChannels = d.Get("slack_channels").([]interface{})
	}
	if d.HasChange("slack_aliases") {
		s.SlackAliases = d.Get("slack_aliases").([]interface{})
	}

	_, err := c.UpdateTeam(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating team: %s", err.Error())
	}

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Team: %s", d.Id()))

	err := c.DeleteTeam(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting team: %s", err.Error())
	}

	d.SetId("")

	return nil
}
