package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceWorkflowTaskCreateJiraSubtask() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow create_jira_subtask task.",

		CreateContext: resourceWorkflowTaskCreateJiraSubtaskCreate,
		ReadContext:   resourceWorkflowTaskCreateJiraSubtaskRead,
		UpdateContext: resourceWorkflowTaskCreateJiraSubtaskUpdate,
		DeleteContext: resourceWorkflowTaskCreateJiraSubtaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Description:  "The ID of the parent workflow",
				Type:         schema.TypeString,
				Required:     true,
			},
			"task_params": {
				Description: "The parameters for this workflow task.",
				Type: schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_type": &schema.Schema{
							Type: schema.TypeString,
							Optional: true,
							Default: "create_jira_subtask",
							ValidateFunc: validation.StringInSlice([]string{
								"create_jira_subtask",
							}, false),
						},
						"project_key": &schema.Schema{
							Description: "The project key.",
							Type: schema.TypeString,
							Required: true,
						},
						"parent_issue_id": &schema.Schema{
							Description: "The parent issue.",
							Type: schema.TypeString,
							Required: true,
						},
						"title": &schema.Schema{
							Description: "The issue title.",
							Type: schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Description: "The issue description.",
							Type: schema.TypeString,
							Optional: true,
						},
						"subtask_issue_type": &schema.Schema{
							Description: "The issue type id and display name.",
							Type: schema.TypeMap,
							Required: true,
						},
						"labels": &schema.Schema{
							Description: "The issue labels.",
							Type: schema.TypeString,
							Optional: true,
						},
						"due_date": &schema.Schema{
							Description: "The due date.",
							Type: schema.TypeString,
							Optional: true,
						},
						"assign_user_email": &schema.Schema{
							Description: "The assigned user's email.",
							Type: schema.TypeString,
							Optional: true,
						},
						"reporter_user_email": &schema.Schema{
							Description: "The reporter user's email.",
							Type: schema.TypeString,
							Optional: true,
						},
						"priority": &schema.Schema{
							Description: "The priority id and display name.",
							Type: schema.TypeMap,
							Optional: true,
						},
						"status": &schema.Schema{
							Description: "The status id and display name.",
							Type: schema.TypeMap,
							Optional: true,
						},
						"custom_fields_mapping": &schema.Schema{
							Description: "",
							Type: schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowTaskCreateJiraSubtaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	workflowId := d.Get("workflow_id").(string)
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow task: %s", workflowId))

	s := &client.WorkflowTask{
		WorkflowId: workflowId,
		TaskParams: taskParams,
	}

	res, err := c.CreateWorkflowTask(s)
	if err != nil {
		return diag.Errorf("Error creating workflow task: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an workflow task resource: %v (%s)", workflowId, d.Id()))

	return resourceWorkflowTaskCreateJiraSubtaskRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateJiraSubtaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow task: %s", d.Id()))

	res, err := c.GetWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateJiraSubtask (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow task: %s", d.Id())
	}

	d.Set("workflow_id", res.WorkflowId)
	tps := make([]interface{}, 1, 1)
	tps[0] = res.TaskParams
	d.Set("task_params", tps)

	return nil
}

func resourceWorkflowTaskCreateJiraSubtaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow task: %s", d.Id()))

	workflowId := d.Get("workflow_id").(string)
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	s := &client.WorkflowTask{
		WorkflowId: workflowId,
		TaskParams: taskParams,
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateWorkflowTask(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow task: %s", err.Error())
	}

	return resourceWorkflowTaskCreateJiraSubtaskRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateJiraSubtaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow task: %s", d.Id()))

	err := c.DeleteWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateJiraSubtask (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow task: %s", err.Error())
	}

	d.SetId("")

	return nil
}