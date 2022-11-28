package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceIncidentRoleTask() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceIncidentRoleTaskCreate,
		ReadContext: resourceIncidentRoleTaskRead,
		UpdateContext: resourceIncidentRoleTaskUpdate,
		DeleteContext: resourceIncidentRoleTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"incident_role_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: true,
				Description: "",
			},
			

			"task": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The task of the incident task",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of incident task",
			},
			

			"priority": &schema.Schema{
				Type: schema.TypeString,
				Default: "high",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The priority of the incident task. Value must be one of `high`, `medium`, `low`.",
			},
			
		},
	}
}


func resourceIncidentRoleTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentRoleTask"))

	s := &client.IncidentRoleTask{}

	  if value, ok := d.GetOkExists("incident_role_id"); ok {
				s.IncidentRoleId = value.(string)
			}
    if value, ok := d.GetOkExists("task"); ok {
				s.Task = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("priority"); ok {
				s.Priority = value.(string)
			}

	res, err := c.CreateIncidentRoleTask(s)
	if err != nil {
		return diag.Errorf("Error creating incident_role_task: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_role_task resource: %s", d.Id()))

	return resourceIncidentRoleTaskRead(ctx, d, meta)
}


func resourceIncidentRoleTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentRoleTask: %s", d.Id()))

	item, err := c.GetIncidentRoleTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRoleTask (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_role_task: %s", d.Id())
	}

	d.Set("incident_role_id", item.IncidentRoleId)
  d.Set("task", item.Task)
  d.Set("description", item.Description)
  d.Set("priority", item.Priority)

	return nil
}


func resourceIncidentRoleTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentRoleTask: %s", d.Id()))

	s := &client.IncidentRoleTask{}

	  if d.HasChange("incident_role_id") {
				s.IncidentRoleId = d.Get("incident_role_id").(string)
			}
    if d.HasChange("task") {
				s.Task = d.Get("task").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}
    if d.HasChange("priority") {
				s.Priority = d.Get("priority").(string)
			}

	_, err := c.UpdateIncidentRoleTask(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_role_task: %s", err.Error())
	}

	return resourceIncidentRoleTaskRead(ctx, d, meta)
}


func resourceIncidentRoleTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentRoleTask: %s", d.Id()))

	err := c.DeleteIncidentRoleTask(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentRoleTask (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_role_task: %s", err.Error())
	}

	d.SetId("")

	return nil
}

