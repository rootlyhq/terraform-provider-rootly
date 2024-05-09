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

func resourceWorkflowActionItem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowActionItemCreate,
		ReadContext:   resourceWorkflowActionItemRead,
		UpdateContext: resourceWorkflowActionItemUpdate,
		DeleteContext: resourceWorkflowActionItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The title of the workflow",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "The slug of the workflow",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "The description of the workflow",
			},

			"command": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Workflow command",
			},

			"command_feedback_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "This will notify you back when the workflow is starting. Value must be one of true or false",
			},

			"wait": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Wait this duration before executing",
			},

			"repeat_every_duration": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Repeat workflow every duration",
			},

			"repeat_on": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "Repeat on weekdays. Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.",
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "The order which the workflow should run with other workflows.",
			},

			"workflow_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "The group this workflow belongs to.",
			},

			"trigger_params": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"trigger_type": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "action_item",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `action_item`.",
						},

						"triggers": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Actions that trigger the workflow. One of custom_fields.<slug>.updated, incident_updated, action_item_created, action_item_updated, assigned_user_updated, summary_updated, description_updated, status_updated, priority_updated, due_date_updated, teams_updated, slack_command",
						},

						"incident_visibilities": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "",
						},

						"incident_kinds": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Value must be one of `test`, `test_sub`, `example`, `example_sub`, `normal`, `normal_sub`, `backfilled`, `scheduled`.",
						},

						"incident_statuses": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Value must be one of `in_triage`, `started`, `detected`, `acknowledged`, `mitigated`, `resolved`, `cancelled`, `scheduled`, `in_progress`, `completed`.",
						},

						"incident_inactivity_duration": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "ex. 10 min, 1h, 3 days, 2 weeks",
						},

						"incident_condition": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ALL",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `ALL`, `ANY`, `NONE`.",
						},

						"incident_condition_visibility": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_kind": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "IS",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_status": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_environment": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_severity": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_incident_type": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_incident_roles": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_service": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_functionality": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_group": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_condition_summary": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_condition_started_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_condition_detected_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_condition_acknowledged_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_condition_mitigated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_condition_resolved_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_conditional_inactivity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"incident_action_item_condition": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ALL",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `ALL`, `ANY`, `NONE`.",
						},

						"incident_action_item_condition_kind": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_action_item_kinds": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Value must be one of `task`, `follow_up`.",
						},

						"incident_action_item_condition_status": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_action_item_statuses": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Value must be one of `open`, `in_progress`, `cancelled`, `done`.",
						},

						"incident_action_item_condition_priority": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_action_item_priorities": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Value must be one of `high`, `medium`, `low`.",
						},

						"incident_action_item_condition_group": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"incident_action_item_group_ids": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "",
						},
					},
				},
				Computed: true,
				Optional: true,
			},

			"environment_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"severity_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"incident_type_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"incident_role_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"service_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"functionality_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"cause_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},
		},
	}
}

func resourceWorkflowActionItemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating WorkflowActionItem"))

	s := &client.Workflow{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("command"); ok {
		s.Command = value.(string)
	}
	if value, ok := d.GetOkExists("command_feedback_enabled"); ok {
		s.CommandFeedbackEnabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("wait"); ok {
		s.Wait = value.(string)
	}
	if value, ok := d.GetOkExists("repeat_every_duration"); ok {
		s.RepeatEveryDuration = value.(string)
	}
	if value, ok := d.GetOkExists("repeat_on"); ok {
		s.RepeatOn = value.([]interface{})
	}
	if value, ok := d.GetOkExists("enabled"); ok {
		s.Enabled = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}
	if value, ok := d.GetOkExists("workflow_group_id"); ok {
		s.WorkflowGroupId = value.(string)
	}
	if value, ok := d.GetOkExists("trigger_params"); ok {
		s.TriggerParams = value.([]interface{})[0].(map[string]interface{})
	}
	if value, ok := d.GetOkExists("environment_ids"); ok {
		s.EnvironmentIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("severity_ids"); ok {
		s.SeverityIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_type_ids"); ok {
		s.IncidentTypeIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("incident_role_ids"); ok {
		s.IncidentRoleIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("service_ids"); ok {
		s.ServiceIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("functionality_ids"); ok {
		s.FunctionalityIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("group_ids"); ok {
		s.GroupIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("cause_ids"); ok {
		s.CauseIds = value.([]interface{})
	}

	res, err := c.CreateWorkflow(s)
	if err != nil {
		return diag.Errorf("Error creating workflow_action_item: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a workflow_action_item resource: %s", d.Id()))

	return resourceWorkflowActionItemRead(ctx, d, meta)
}

func resourceWorkflowActionItemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WorkflowActionItem: %s", d.Id()))

	item, err := c.GetWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowActionItem (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow_action_item: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("command", item.Command)
	d.Set("command_feedback_enabled", item.CommandFeedbackEnabled)
	d.Set("wait", item.Wait)
	d.Set("repeat_every_duration", item.RepeatEveryDuration)
	d.Set("repeat_on", item.RepeatOn)
	d.Set("enabled", item.Enabled)
	d.Set("position", item.Position)
	d.Set("workflow_group_id", item.WorkflowGroupId)

	tps := make([]interface{}, 1, 1)
	tps[0] = item.TriggerParams
	d.Set("trigger_params", tps)

	d.Set("environment_ids", item.EnvironmentIds)
	d.Set("severity_ids", item.SeverityIds)
	d.Set("incident_type_ids", item.IncidentTypeIds)
	d.Set("incident_role_ids", item.IncidentRoleIds)
	d.Set("service_ids", item.ServiceIds)
	d.Set("functionality_ids", item.FunctionalityIds)
	d.Set("group_ids", item.GroupIds)
	d.Set("cause_ids", item.CauseIds)

	return nil
}

func resourceWorkflowActionItemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating WorkflowActionItem: %s", d.Id()))

	s := &client.Workflow{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("command") {
		s.Command = d.Get("command").(string)
	}
	if d.HasChange("command_feedback_enabled") {
		s.CommandFeedbackEnabled = tools.Bool(d.Get("command_feedback_enabled").(bool))
	}
	if d.HasChange("wait") {
		s.Wait = d.Get("wait").(string)
	}
	if d.HasChange("repeat_every_duration") {
		s.RepeatEveryDuration = d.Get("repeat_every_duration").(string)
	}
	if d.HasChange("repeat_on") {
		s.RepeatOn = d.Get("repeat_on").([]interface{})
	}
	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}
	if d.HasChange("workflow_group_id") {
		s.WorkflowGroupId = d.Get("workflow_group_id").(string)
	}

	if d.HasChange("trigger_params") {
		tps := d.Get("trigger_params").([]interface{})
		for _, tpsi := range tps {
			s.TriggerParams = tpsi.(map[string]interface{})
		}
	}
	if d.HasChange("environment_ids") {
		s.EnvironmentIds = d.Get("environment_ids").([]interface{})
	}
	if d.HasChange("severity_ids") {
		s.SeverityIds = d.Get("severity_ids").([]interface{})
	}
	if d.HasChange("incident_type_ids") {
		s.IncidentTypeIds = d.Get("incident_type_ids").([]interface{})
	}
	if d.HasChange("incident_role_ids") {
		s.IncidentRoleIds = d.Get("incident_role_ids").([]interface{})
	}
	if d.HasChange("service_ids") {
		s.ServiceIds = d.Get("service_ids").([]interface{})
	}
	if d.HasChange("functionality_ids") {
		s.FunctionalityIds = d.Get("functionality_ids").([]interface{})
	}
	if d.HasChange("group_ids") {
		s.GroupIds = d.Get("group_ids").([]interface{})
	}
	if d.HasChange("cause_ids") {
		s.CauseIds = d.Get("cause_ids").([]interface{})
	}

	_, err := c.UpdateWorkflow(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow_action_item: %s", err.Error())
	}

	return resourceWorkflowActionItemRead(ctx, d, meta)
}

func resourceWorkflowActionItemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting WorkflowActionItem: %s", d.Id()))

	err := c.DeleteWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowActionItem (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow_action_item: %s", err.Error())
	}

	d.SetId("")

	return nil
}
