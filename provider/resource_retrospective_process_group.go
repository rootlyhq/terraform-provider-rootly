package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceRetrospectiveProcessGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRetrospectiveProcessGroupCreate,
		ReadContext:   resourceRetrospectiveProcessGroupRead,
		UpdateContext: resourceRetrospectiveProcessGroupUpdate,
		DeleteContext: resourceRetrospectiveProcessGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"retrospective_process_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},

			"sub_status_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},
		},
	}
}

func resourceRetrospectiveProcessGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating RetrospectiveProcessGroup"))

	s := &client.RetrospectiveProcessGroup{}

	if value, ok := d.GetOkExists("retrospective_process_id"); ok {
		s.RetrospectiveProcessId = value.(string)
	}
	if value, ok := d.GetOkExists("sub_status_id"); ok {
		s.SubStatusId = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateRetrospectiveProcessGroup(s)
	if err != nil {
		return diag.Errorf("Error creating retrospective_process_group: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a retrospective_process_group resource: %s", d.Id()))

	return resourceRetrospectiveProcessGroupRead(ctx, d, meta)
}

func resourceRetrospectiveProcessGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading RetrospectiveProcessGroup: %s", d.Id()))

	item, err := c.GetRetrospectiveProcessGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcessGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading retrospective_process_group: %s", d.Id())
	}

	d.Set("retrospective_process_id", item.RetrospectiveProcessId)
	d.Set("sub_status_id", item.SubStatusId)
	d.Set("position", item.Position)

	return nil
}

func resourceRetrospectiveProcessGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating RetrospectiveProcessGroup: %s", d.Id()))

	s := &client.RetrospectiveProcessGroup{}

	if d.HasChange("retrospective_process_id") {
		s.RetrospectiveProcessId = d.Get("retrospective_process_id").(string)
	}
	if d.HasChange("sub_status_id") {
		s.SubStatusId = d.Get("sub_status_id").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateRetrospectiveProcessGroup(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating retrospective_process_group: %s", err.Error())
	}

	return resourceRetrospectiveProcessGroupRead(ctx, d, meta)
}

func resourceRetrospectiveProcessGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting RetrospectiveProcessGroup: %s", d.Id()))

	err := c.DeleteRetrospectiveProcessGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcessGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting retrospective_process_group: %s", err.Error())
	}

	d.SetId("")

	return nil
}