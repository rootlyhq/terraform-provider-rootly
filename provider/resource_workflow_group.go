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

func resourceWorkflowGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow groups.",

		CreateContext: resourceWorkflowGroupCreate,
		ReadContext:   resourceWorkflowGroupRead,
		UpdateContext: resourceWorkflowGroupUpdate,
		DeleteContext: resourceWorkflowGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the workflow group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kind": {
				Description: "The kind of the workflow group",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"incident",
					"post_mortem",
					"action_item",
					"alert",
					"pulse",
					"simple",
				}, false),
			},
			"expanded": {
				Description: "Whether the workflow group is expanded or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"position": {
				Description:  "The position of the workflow group (1 being top of list)",
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
			},
		},
	}
}

func resourceWorkflowGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	kind := d.Get("kind").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow group: %s", name))

	s := &client.WorkflowGroup{
		Name: name,
		Kind: kind,
	}

	if v, ok := d.GetOkExists("expanded"); ok {
		s.Expanded = tools.Bool(v.(bool))
	}

	if v, ok := d.GetOk("position"); ok {
		s.Position = v.(int)
	}

	res, err := c.CreateWorkflowGroup(s)
	if err != nil {
		return diag.Errorf("Error creating workflow group: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an workflow group resource: %v (%s)", name, d.Id()))

	return resourceWorkflowGroupRead(ctx, d, meta)
}

func resourceWorkflowGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow group: %s", d.Id()))

	res, err := c.GetWorkflowGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow group: %s", d.Id())
	}

	d.Set("name", res.Name)
	d.Set("kind", res.Kind)
	d.Set("expanded", res.Expanded)
	d.Set("position", res.Position)

	return nil
}

func resourceWorkflowGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow group: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.WorkflowGroup{
		Name: name,
	}

	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}

	if d.HasChange("expanded") {
		s.Expanded = tools.Bool(d.Get("expanded").(bool))
	}

	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	tflog.Debug(ctx, fmt.Sprintf("adding value: %#v", s))
	_, err := c.UpdateWorkflowGroup(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow group: %s", err.Error())
	}

	return resourceWorkflowGroupRead(ctx, d, meta)
}

func resourceWorkflowGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow group: %s", d.Id()))

	err := c.DeleteWorkflowGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow group: %s", err.Error())
	}

	d.SetId("")

	return nil
}
