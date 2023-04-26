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

func resourceWorkflowAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowAlertCreate,
		ReadContext:   resourceWorkflowAlertRead,
		UpdateContext: resourceWorkflowAlertUpdate,
		DeleteContext: resourceWorkflowAlertDelete,
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
				Description: "Workflow command.",
			},

			"wait": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Wait this duration before executing.",
			},

			"repeat_every_duration": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Repeat workflow every duration.",
			},

			"repeat_on": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.",
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
							Default:     "alert",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `alert`.",
						},

						"triggers": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "Value must be one of `alert_created`.",
						},

						"alert_condition": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ALL",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `ALL`, `ANY`, `NONE`.",
						},

						"alert_condition_source": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"alert_condition_source_use_regexp": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_sources": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_condition_label": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"alert_condition_label_use_regexp": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_labels": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_condition_payload": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "ANY",
							Required:    false,
							Optional:    true,
							Description: "Value must be one off `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
						},

						"alert_condition_payload_use_regexp": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_payload": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "",
						},

						"alert_query_payload": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "You can use jsonpath syntax. eg: $.incident.teams[*]",
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
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},

			"severity_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},

			"incident_type_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},

			"incident_role_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},

			"service_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},

			"group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceWorkflowAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating WorkflowAlert"))

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
	if value, ok := d.GetOkExists("group_ids"); ok {
		s.GroupIds = value.([]interface{})
	}

	res, err := c.CreateWorkflow(s)
	if err != nil {
		return diag.Errorf("Error creating workflow_alert: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a workflow_alert resource: %s", d.Id()))

	return resourceWorkflowAlertRead(ctx, d, meta)
}

func resourceWorkflowAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WorkflowAlert: %s", d.Id()))

	item, err := c.GetWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowAlert (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow_alert: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("command", item.Command)
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
	d.Set("group_ids", item.GroupIds)

	return nil
}

func resourceWorkflowAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating WorkflowAlert: %s", d.Id()))

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
	if d.HasChange("group_ids") {
		s.GroupIds = d.Get("group_ids").([]interface{})
	}

	_, err := c.UpdateWorkflow(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow_alert: %s", err.Error())
	}

	return resourceWorkflowAlertRead(ctx, d, meta)
}

func resourceWorkflowAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting WorkflowAlert: %s", d.Id()))

	err := c.DeleteWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowAlert (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow_alert: %s", err.Error())
	}

	d.SetId("")

	return nil
}
