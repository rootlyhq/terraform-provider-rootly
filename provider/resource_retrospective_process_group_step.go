// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	
)

func resourceRetrospectiveProcessGroupStep() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRetrospectiveProcessGroupStepCreate,
		ReadContext: resourceRetrospectiveProcessGroupStepRead,
		UpdateContext: resourceRetrospectiveProcessGroupStepUpdate,
		DeleteContext: resourceRetrospectiveProcessGroupStepDelete,
		Importer: &schema.ResourceImporter {
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema {
			
			"retrospective_process_group_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: true,
				Description: "",
				
			},
			

			"retrospective_step_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "",
				
			},
			

		"position": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "",
			
		},
		
		},
	}
}

func resourceRetrospectiveProcessGroupStepCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating RetrospectiveProcessGroupStep"))

	s := &client.RetrospectiveProcessGroupStep{}

	  if value, ok := d.GetOkExists("retrospective_process_group_id"); ok {
				s.RetrospectiveProcessGroupId = value.(string)
			}
    if value, ok := d.GetOkExists("retrospective_step_id"); ok {
				s.RetrospectiveStepId = value.(string)
			}
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}

	res, err := c.CreateRetrospectiveProcessGroupStep(s)
	if err != nil {
		return diag.Errorf("Error creating retrospective_process_group_step: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a retrospective_process_group_step resource: %s", d.Id()))

	return resourceRetrospectiveProcessGroupStepRead(ctx, d, meta)
}

func resourceRetrospectiveProcessGroupStepRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading RetrospectiveProcessGroupStep: %s", d.Id()))

	item, err := c.GetRetrospectiveProcessGroupStep(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcessGroupStep (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading retrospective_process_group_step: %s", d.Id())
	}

	d.Set("retrospective_process_group_id", item.RetrospectiveProcessGroupId)
  d.Set("retrospective_step_id", item.RetrospectiveStepId)
  d.Set("position", item.Position)

	return nil
}

func resourceRetrospectiveProcessGroupStepUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating RetrospectiveProcessGroupStep: %s", d.Id()))

	s := &client.RetrospectiveProcessGroupStep{}

	  if d.HasChange("retrospective_process_group_id") {
				s.RetrospectiveProcessGroupId = d.Get("retrospective_process_group_id").(string)
			}
    if d.HasChange("retrospective_step_id") {
				s.RetrospectiveStepId = d.Get("retrospective_step_id").(string)
			}
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}

	_, err := c.UpdateRetrospectiveProcessGroupStep(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating retrospective_process_group_step: %s", err.Error())
	}

	return resourceRetrospectiveProcessGroupStepRead(ctx, d, meta)
}

func resourceRetrospectiveProcessGroupStepDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting RetrospectiveProcessGroupStep: %s", d.Id()))

	err := c.DeleteRetrospectiveProcessGroupStep(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcessGroupStep (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting retrospective_process_group_step: %s", err.Error())
	}

	d.SetId("")

	return nil
}
