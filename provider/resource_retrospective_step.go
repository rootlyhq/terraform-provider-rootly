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

func resourceRetrospectiveStep() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRetrospectiveStepCreate,
		ReadContext:   resourceRetrospectiveStepRead,
		UpdateContext: resourceRetrospectiveStepUpdate,
		DeleteContext: resourceRetrospectiveStepDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"title": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the step",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The slug of the step",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the step",
			},

			"incident_role_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Users assigned to the selected incident role will be the default owners for this step",
			},

			"due_after_days": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Due date in days",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Position of the step",
			},

			"skippable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Is the step skippable?",
			},
		},
	}
}

func resourceRetrospectiveStepCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating RetrospectiveStep"))

	s := &client.RetrospectiveStep{}

	if value, ok := d.GetOkExists("title"); ok {
		s.Title = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("incident_role_id"); ok {
		s.IncidentRoleId = value.(string)
	}
	if value, ok := d.GetOkExists("due_after_days"); ok {
		s.DueAfterDays = value.(int)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}
	if value, ok := d.GetOkExists("skippable"); ok {
		s.Skippable = tools.Bool(value.(bool))
	}

	res, err := c.CreateRetrospectiveStep(s)
	if err != nil {
		return diag.Errorf("Error creating retrospective_step: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a retrospective_step resource: %s", d.Id()))

	return resourceRetrospectiveStepRead(ctx, d, meta)
}

func resourceRetrospectiveStepRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading RetrospectiveStep: %s", d.Id()))

	item, err := c.GetRetrospectiveStep(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveStep (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading retrospective_step: %s", d.Id())
	}

	d.Set("title", item.Title)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("incident_role_id", item.IncidentRoleId)
	d.Set("due_after_days", item.DueAfterDays)
	d.Set("position", item.Position)
	d.Set("skippable", item.Skippable)

	return nil
}

func resourceRetrospectiveStepUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating RetrospectiveStep: %s", d.Id()))

	s := &client.RetrospectiveStep{}

	if d.HasChange("title") {
		s.Title = d.Get("title").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("incident_role_id") {
		s.IncidentRoleId = d.Get("incident_role_id").(string)
	}
	if d.HasChange("due_after_days") {
		s.DueAfterDays = d.Get("due_after_days").(int)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}
	if d.HasChange("skippable") {
		s.Skippable = tools.Bool(d.Get("skippable").(bool))
	}

	_, err := c.UpdateRetrospectiveStep(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating retrospective_step: %s", err.Error())
	}

	return resourceRetrospectiveStepRead(ctx, d, meta)
}

func resourceRetrospectiveStepDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting RetrospectiveStep: %s", d.Id()))

	err := c.DeleteRetrospectiveStep(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveStep (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting retrospective_step: %s", err.Error())
	}

	d.SetId("")

	return nil
}
