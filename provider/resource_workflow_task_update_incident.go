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

func resourceWorkflowTaskUpdateIncident() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow update_incident task.",

		CreateContext: resourceWorkflowTaskUpdateIncidentCreate,
		ReadContext:   resourceWorkflowTaskUpdateIncidentRead,
		UpdateContext: resourceWorkflowTaskUpdateIncidentUpdate,
		DeleteContext: resourceWorkflowTaskUpdateIncidentDelete,
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
							Default: "update_incident",
							ValidateFunc: validation.StringInSlice([]string{
								"update_incident",
							}, false),
						},
						"incident_id": &schema.Schema{
							Description: "The incident id to update",
							Type: schema.TypeString,
							Optional: true,
						},
						"title": &schema.Schema{
							Description: "The incident title",
							Type: schema.TypeString,
							Optional: true,
						},
						"summary": &schema.Schema{
							Description: "The incident summary",
							Type: schema.TypeString,
							Optional: true,
						},
						"status": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"severity_id": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"incident_type_ids": &schema.Schema{
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service_ids": &schema.Schema{
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"functionality_ids": &schema.Schema{
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"environment_ids": &schema.Schema{
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"group_ids": &schema.Schema{
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"started_at": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"detected_at": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"acknowledged_at": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"mitigated_at": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"resolved_at": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"private": &schema.Schema{
							Description: "",
							Type: schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowTaskUpdateIncidentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceWorkflowTaskUpdateIncidentRead(ctx, d, meta)
}

func resourceWorkflowTaskUpdateIncidentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow task: %s", d.Id()))

	res, err := c.GetWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskUpdateIncident (%s) not found, removing from state", d.Id()))
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

func resourceWorkflowTaskUpdateIncidentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceWorkflowTaskUpdateIncidentRead(ctx, d, meta)
}

func resourceWorkflowTaskUpdateIncidentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow task: %s", d.Id()))

	err := c.DeleteWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskUpdateIncident (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow task: %s", err.Error())
	}

	d.SetId("")

	return nil
}