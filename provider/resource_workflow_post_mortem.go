package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	"github.com/rootlyhq/terraform-provider-rootly/tools"
)

func resourceWorkflowPostMortem() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflows.",

		CreateContext: resourceWorkflowPostMortemCreate,
		ReadContext:   resourceWorkflowPostMortemRead,
		UpdateContext: resourceWorkflowPostMortemUpdate,
		DeleteContext: resourceWorkflowPostMortemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the workflow",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the workflow",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the workflow is enabled or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"command": {
				Description: "The workflow command.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"wait": {
				Description: "Wait before running workflow.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"repeat_every_duration": {
				Description: "Repeat workflow every duration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"repeat_on": {
				Description: "Repeat workflow on days.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"S",
						"M",
						"T",
						"W",
						"R",
						"F",
						"U",
					}, false),
				},
			},
			"environment_ids": {
				Description: "Environment IDs required to trigger workflow.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"severity_ids": {
				Description: "Severity IDs required to trigger workflow.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"incident_type_ids": {
				Description: "Incident type IDs required to trigger workflow.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_ids": {
				Description: "Service IDs required to trigger workflow.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_ids": {
				Description: "Group IDs required to trigger workflow.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"trigger_params": {
				Description: "The conditions for triggering this workflow.",
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_type": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "post_mortem",
							ValidateFunc: validation.StringInSlice([]string{
								"post_mortem",
							}, false),
						},
						"triggers": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"incident_kinds": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"incident_statuses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"incident_visibilities": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
						"incident_condition": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ALL",
						},
						"incident_condition_visibility": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_kind": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_status": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_environment": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_severity": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_incident_type": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_service": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_functionality": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_condition_group": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "ANY",
						},
						"incident_post_mortem_condition": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default: "ALL",
						},
						"incident_post_mortem_condition_status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default: "ANY",
						},
						"incident_post_mortem_condition_cause": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default: "ANY",
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowPostMortemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow: %s", name))

	s := &client.Workflow{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	if v, ok := d.GetOkExists("enabled"); ok {
		s.Enabled = tools.Bool(v.(bool))
	}

	if value, ok := d.GetOk("command"); ok {
		s.Command = value.(string)
	}

	if value, ok := d.GetOk("wait"); ok {
		s.Wait = value.(string)
	}

	if value, ok := d.GetOk("trigger_params"); ok {
		s.TriggerParams = value.([]interface{})[0].(map[string]interface{})
	}

	if value, ok := d.GetOk("repeat_every_duration"); ok {
		s.RepeatEveryDuration = value.(string)
	}

	if value, ok := d.GetOk("repeat_on"); ok {
		s.RepeatOn = value.([]interface{})
	}

	if value, ok := d.GetOk("environment_ids"); ok {
		s.EnvironmentIds = value.([]interface{})
	}

	if value, ok := d.GetOk("severity_ids"); ok {
		s.SeverityIds = value.([]interface{})
	}

	if value, ok := d.GetOk("incident_type_ids"); ok {
		s.IncidentTypeIds = value.([]interface{})
	}

	if value, ok := d.GetOk("service_ids"); ok {
		s.ServiceIds = value.([]interface{})
	}

	if value, ok := d.GetOk("group_ids"); ok {
		s.GroupIds = value.([]interface{})
	}

	res, err := c.CreateWorkflow(s)
	if err != nil {
		return diag.Errorf("Error creating workflow: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an workflow resource: %v (%s)", name, d.Id()))

	return resourceWorkflowPostMortemRead(ctx, d, meta)
}

func resourceWorkflowPostMortemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow: %s", d.Id()))

	res, err := c.GetWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Workflow (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow: %s", d.Id())
	}

	d.Set("name", res.Name)
	d.Set("description", res.Description)
	d.Set("enabled", res.Enabled)
	tps := make([]interface{}, 1, 1)
	tps[0] = res.TriggerParams
	d.Set("trigger_params", tps)
	d.Set("command", res.Command)
	d.Set("wait", res.Wait)
	d.Set("repeat_every_duration", res.RepeatEveryDuration)
	d.Set("repeat_on", res.RepeatOn)
	d.Set("environment_ids", res.EnvironmentIds)
	d.Set("severity_ids", res.SeverityIds)
	d.Set("incident_type_ids", res.IncidentTypeIds)
	d.Set("service_ids", res.ServiceIds)
	d.Set("group_ids", res.GroupIds)

	return nil
}

func resourceWorkflowPostMortemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Workflow{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("enabled") {
		s.Enabled = tools.Bool(d.Get("enabled").(bool))
	}

	if d.HasChange("trigger_params") {
		tps := d.Get("trigger_params").([]interface{})
		for _, tpsi := range tps {
			s.TriggerParams = tpsi.(map[string]interface{})
		}
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

	if d.HasChange("environment_ids") {
		s.EnvironmentIds = d.Get("environment_ids").([]interface{})
	}

	if d.HasChange("severity_ids") {
		s.SeverityIds = d.Get("severity_ids").([]interface{})
	}

	if d.HasChange("incident_type_ids") {
		s.IncidentTypeIds = d.Get("incident_type_ids").([]interface{})
	}

	if d.HasChange("service_ids") {
		s.ServiceIds = d.Get("service_ids").([]interface{})
	}

	if d.HasChange("group_ids") {
		s.GroupIds = d.Get("group_ids").([]interface{})
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateWorkflow(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow: %s", err.Error())
	}

	return resourceWorkflowPostMortemRead(ctx, d, meta)
}

func resourceWorkflowPostMortemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow: %s", d.Id()))

	err := c.DeleteWorkflow(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Workflow (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow: %s", err.Error())
	}

	d.SetId("")

	return nil
}
