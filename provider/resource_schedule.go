package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

// Upgrade logic from version 0 to version 1
func upgradeScheduleStateV0ToV1(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	delete(rawState, "slack_user_group")
	return rawState, nil
}

func resourceScheduleV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"slack_user_group": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Computed: true,
			},
		},
	}
}

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
				Type:    resourceScheduleV0().CoreConfigSchema().ImpliedType(),
				Upgrade: upgradeScheduleStateV0ToV1,
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

			"owner_group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
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
				Description: "ID of user assigned as owner of the schedule",
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
	if value, ok := d.GetOkExists("owner_group_ids"); ok {
		s.OwnerGroupIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("owner_user_id"); ok {
		s.OwnerUserId = value.(int)
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
	d.Set("owner_group_ids", item.OwnerGroupIds)
	d.Set("owner_user_id", item.OwnerUserId)

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
	if d.HasChange("owner_group_ids") {
		s.OwnerGroupIds = d.Get("owner_group_ids").([]interface{})
	}
	if d.HasChange("owner_user_id") {
		s.OwnerUserId = d.Get("owner_user_id").(int)
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
