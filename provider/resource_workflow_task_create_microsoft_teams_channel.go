package provider

// This file was auto-generated by tools/generate-tasks.js

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceWorkflowTaskCreateMicrosoftTeamsChannel() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow create_microsoft_teams_channel task.",

		CreateContext: resourceWorkflowTaskCreateMicrosoftTeamsChannelCreate,
		ReadContext:   resourceWorkflowTaskCreateMicrosoftTeamsChannelRead,
		UpdateContext: resourceWorkflowTaskCreateMicrosoftTeamsChannelUpdate,
		DeleteContext: resourceWorkflowTaskCreateMicrosoftTeamsChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Description: "The ID of the parent workflow",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the workflow task",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"position": {
				Description: "The position of the workflow task (1 being top of list)",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"skip_on_failure": {
				Description: "Skip workflow task if any failures",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"enabled": {
				Description: "Enable/disable this workflow task",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"task_params": {
				Description: "The parameters for this workflow task.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "create_microsoft_teams_channel",
							ValidateFunc: validation.StringInSlice([]string{
								"create_microsoft_teams_channel",
							}, false),
						},
						"team": &schema.Schema{
							Description: "Map must contain two fields, `id` and `name`. ",
							Type:        schema.TypeMap,
							Optional:    true,
						},
						"title": &schema.Schema{
							Description: "Microsoft Team channel title",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": &schema.Schema{
							Description: "Microsoft Team channel description",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"private": &schema.Schema{
							Description: "Value must be one of `auto`, `true`, `false`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "auto",
							ValidateFunc: validation.StringInSlice([]string{
								"auto",
								"true",
								"false",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowTaskCreateMicrosoftTeamsChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	workflowId := d.Get("workflow_id").(string)
	name := d.Get("name").(string)
	position := d.Get("position").(int)
	skipOnFailure := tools.Bool(d.Get("skip_on_failure").(bool))
	enabled := tools.Bool(d.Get("enabled").(bool))
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow task: %s", workflowId))

	s := &client.WorkflowTask{
		WorkflowId:    workflowId,
		Name:          name,
		Position:      position,
		SkipOnFailure: skipOnFailure,
		Enabled:       enabled,
		TaskParams:    taskParams,
	}

	res, err := c.CreateWorkflowTask(s)
	if err != nil {
		return diag.Errorf("Error creating workflow task: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an workflow task resource: %v (%s)", workflowId, d.Id()))

	return resourceWorkflowTaskCreateMicrosoftTeamsChannelRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateMicrosoftTeamsChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow task: %s", d.Id()))

	res, err := c.GetWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateMicrosoftTeamsChannel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow task: %s", d.Id())
	}

	d.Set("workflow_id", res.WorkflowId)
	d.Set("name", res.Name)
	d.Set("position", res.Position)
	d.Set("skip_on_failure", res.SkipOnFailure)
	d.Set("enabled", res.Enabled)
	tps := make([]interface{}, 1, 1)
	tps[0] = res.TaskParams
	d.Set("task_params", tps)

	return nil
}

func resourceWorkflowTaskCreateMicrosoftTeamsChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow task: %s", d.Id()))

	workflowId := d.Get("workflow_id").(string)
	name := d.Get("name").(string)
	position := d.Get("position").(int)
	skipOnFailure := tools.Bool(d.Get("skip_on_failure").(bool))
	enabled := tools.Bool(d.Get("enabled").(bool))
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	s := &client.WorkflowTask{
		WorkflowId:    workflowId,
		Name:          name,
		Position:      position,
		SkipOnFailure: skipOnFailure,
		Enabled:       enabled,
		TaskParams:    taskParams,
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateWorkflowTask(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow task: %s", err.Error())
	}

	return resourceWorkflowTaskCreateMicrosoftTeamsChannelRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateMicrosoftTeamsChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow task: %s", d.Id()))

	err := c.DeleteWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateMicrosoftTeamsChannel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow task: %s", err.Error())
	}

	d.SetId("")

	return nil
}