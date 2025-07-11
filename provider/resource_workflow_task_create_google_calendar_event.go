package provider

// This file was auto-generated by tools/generate-tasks.js

import (
	"context"
	"errors"
	"fmt"
	
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
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

		Schema: map[string]*schema.Schema {
			"workflow_id": {
				Description:  "The ID of the parent workflow",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"name": {
				Description:  "Name of the workflow task",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
			},
			"position": {
				Description:  "The position of the workflow task (1 being top of list)",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
			},
			"skip_on_failure": {
				Description:  "Skip workflow task if any failures",
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
			},
			"enabled": {
				Description:  "Enable/disable this workflow task",
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
			},
			"task_params": {
				Description: "The parameters for this workflow task.",
				Type: schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema {
						"task_type": &schema.Schema {
							Type: schema.TypeString,
							Optional: true,
							Default: "create_google_calendar_event",
							ValidateFunc: validation.StringInSlice([]string {
								"create_google_calendar_event",
							}, false),
						},
						"attendees": &schema.Schema {
							Description: "Emails of attendees",
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Schema {
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
						},
						"time_zone": &schema.Schema {
							Description: "A valid IANA time zone name.",
							Type: schema.TypeString,
							Optional: true,
						},
						"calendar_id": &schema.Schema {
							Description: "",
							Type: schema.TypeString,
							Optional: true,
						},
						"days_until_meeting": &schema.Schema {
							Description: "The days until meeting",
							Type: schema.TypeString,
							Required: true,
						},
						"time_of_meeting": &schema.Schema {
							Description: "Time of meeting in format HH:MM",
							Type: schema.TypeString,
							Required: true,
						},
						"meeting_duration": &schema.Schema {
							Description: "Meeting duration in format like '1 hour', '30 minutes'",
							Type: schema.TypeString,
							Required: true,
						},
						"send_updates": &schema.Schema {
							Description: "Send an email to the attendees notifying them of the event. Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_modify_event": &schema.Schema {
							Description: "Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_see_other_guests": &schema.Schema {
							Description: "Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"can_guests_invite_others": &schema.Schema {
							Description: "Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"summary": &schema.Schema {
							Description: "The event summary",
							Type: schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema {
							Description: "The event description",
							Type: schema.TypeString,
							Required: true,
						},
						"exclude_weekends": &schema.Schema {
							Description: "Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"conference_solution_key": &schema.Schema {
							Description: "Sets the video conference type attached to the meeting. Value must be one of `eventHangout`, `eventNamedHangout`, `hangoutsMeet`, `addOn`.",
							Type: schema.TypeString,
							Optional: true,
							Default: nil,
							ValidateFunc: validation.StringInSlice([]string{
								"eventHangout",
"eventNamedHangout",
"hangoutsMeet",
"addOn",
							}, false),
						},
						"post_to_incident_timeline": &schema.Schema {
							Description: "Value must be one of true or false",
							Type: schema.TypeBool,
							Optional: true,
						},
						"post_to_slack_channels": &schema.Schema {
							Description: "",
							Type: schema.TypeList,
							Optional: true,
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema {
									"id": &schema.Schema {
										Type: schema.TypeString,
										Required: true,
									},
									"name": &schema.Schema {
										Type: schema.TypeString,
										Required: true,
									},
								},
							},
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
	name := d.Get("name").(string)
	position := d.Get("position").(int)
	skipOnFailure := tools.Bool(d.Get("skip_on_failure").(bool))
	enabled := tools.Bool(d.Get("enabled").(bool))
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow task: %s", workflowId))

	s := &client.WorkflowTask{
		WorkflowId: workflowId,
		Name: name,
		Position: position,
		SkipOnFailure: skipOnFailure,
		Enabled: enabled,
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
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateGoogleCalendarEvent (%s) not found, removing from state", d.Id()))
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

func resourceWorkflowTaskCreateGoogleCalendarEventUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow task: %s", d.Id()))

	workflowId := d.Get("workflow_id").(string)
	name := d.Get("name").(string)
	position := d.Get("position").(int)
	skipOnFailure := tools.Bool(d.Get("skip_on_failure").(bool))
	enabled := tools.Bool(d.Get("enabled").(bool))
	taskParams := d.Get("task_params").([]interface{})[0].(map[string]interface{})

	s := &client.WorkflowTask{
		WorkflowId: workflowId,
		Name: name,
		Position: position,
		SkipOnFailure: skipOnFailure,
		Enabled: enabled,
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
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowTaskCreateGoogleCalendarEvent (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow task: %s", err.Error())
	}

	d.SetId("")

	return nil
}
