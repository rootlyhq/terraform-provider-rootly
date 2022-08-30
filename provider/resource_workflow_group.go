package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceWorkflowGroup() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceWorkflowGroupCreate,
		ReadContext: resourceWorkflowGroupRead,
		UpdateContext: resourceWorkflowGroupUpdate,
		DeleteContext: resourceWorkflowGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"kind": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The kind of the workflow group.",
			},
			

			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The name of the workflow group.",
			},
			

			"expanded": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "Whether the group is expanded or collapsed.",
			},
			

			"position": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The position of the workflow group",
			},
			
		},
	}
}

func resourceWorkflowGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating WorkflowGroup"))

	s := &client.WorkflowGroup{}

	  if value, ok := d.GetOkExists("kind"); ok {
		s.Kind = value.(string)
	}
    if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
    if value, ok := d.GetOkExists("expanded"); ok {
		s.Expanded = value.(bool)
	}
    if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateWorkflowGroup(s)
	if err != nil {
		return diag.Errorf("Error creating workflow_group: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a workflow_group resource: %s", d.Id()))

	return resourceWorkflowGroupRead(ctx, d, meta)
}

func resourceWorkflowGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WorkflowGroup: %s", d.Id()))

	item, err := c.GetWorkflowGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow_group: %s", d.Id())
	}

	d.Set("kind", item.Kind)
  d.Set("name", item.Name)
  d.Set("expanded", item.Expanded)
  d.Set("position", item.Position)

	return nil
}

func resourceWorkflowGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating WorkflowGroup: %s", d.Id()))

	s := &client.WorkflowGroup{}

	  if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
    if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
    if d.HasChange("expanded") {
		s.Expanded = d.Get("expanded").(bool)
	}
    if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateWorkflowGroup(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow_group: %s", err.Error())
	}

	return resourceWorkflowGroupRead(ctx, d, meta)
}

func resourceWorkflowGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting WorkflowGroup: %s", d.Id()))

	err := c.DeleteWorkflowGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow_group: %s", err.Error())
	}

	d.SetId("")

	return nil
}
