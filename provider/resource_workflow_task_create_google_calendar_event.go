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

func resourceWorkflowTaskCreateGoogleCalendarEvent() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow create_google_calendar_event task.",

		CreateContext: resourceWorkflowTaskCreateGoogleCalendarEventCreate,
		ReadContext:   resourceWorkflowTaskCreateGoogleCalendarEventRead,
		UpdateContext: resourceWorkflowTaskCreateGoogleCalendarEventUpdate,
		DeleteContext: resourceWorkflowTaskCreateGoogleCalendarEventDelete,
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
							Default: "create_google_calendar_event",
							ValidateFunc: validation.StringInSlice([]string{
								"create_google_calendar_event",
							}, false),
						},
						"attendees": &schema.Schema{
							Description: "Emails of attendees",
							Type: schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"time_zone": &schema.Schema{
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"days_until_meeting": &schema.Schema{
							Description: "The days until meeting",
							Type: schema.TypeInt,
							Required: true,
						},
						"time_of_meeting": &schema.Schema{
							Description: "Time of meeting in format HH:MM",
							Type: schema.TypeString,
							Required: true,
						},
						"meeting_duration": &schema.Schema{
							Description: "Meeting duration in format like '1 hour', '30 minutes'",
							Type: schema.TypeString,
							Required: true,
						},
						"send_updates": &schema.Schema{
							Description: "Send an email to the attendees notifying them of the event",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_modify_event": &schema.Schema{
							Description: "",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_see_other_guests": &schema.Schema{
							Description: "",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_invite_others": &schema.Schema{
							Description: "",
							Type: schema.TypeBool,
							Optional: true,
						},
						"post_to_incident_timeline": &schema.Schema{
							Description: "",
							Type: schema.TypeBool,
							Optional: true,
						},
						"summary": &schema.Schema{
							Description: "The event summary",
							Type: schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Description: "The event description",
							Type: schema.TypeString,
							Required: true,
						},
						"exclude_weekends": &schema.Schema{
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

func resourceWorkflowTaskCreateGoogleCalendarEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceWorkflowTaskCreateGoogleCalendarEventRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateGoogleCalendarEventRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow task: %s", d.Id()))

	res, err := c.GetWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateGoogleCalendarEvent (%s) not found, removing from state", d.Id()))
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

func resourceWorkflowTaskCreateGoogleCalendarEventUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceWorkflowTaskCreateGoogleCalendarEventRead(ctx, d, meta)
}

func resourceWorkflowTaskCreateGoogleCalendarEventDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow task: %s", d.Id()))

	err := c.DeleteWorkflowTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateGoogleCalendarEvent (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow task: %s", err.Error())
	}

	d.SetId("")

	return nil
}
