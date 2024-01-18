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

func resourceRetrospectiveProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRetrospectiveProcessCreate,
		ReadContext:   resourceRetrospectiveProcessRead,
		UpdateContext: resourceRetrospectiveProcessUpdate,
		DeleteContext: resourceRetrospectiveProcessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the retrospective process",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the retrospective process",
			},

			"is_default": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Is the retrospective process default?. Value must be one of true or false",
			},

			"retrospective_process_matching_criteria": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceRetrospectiveProcessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating RetrospectiveProcess"))

	s := &client.RetrospectiveProcess{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("is_default"); ok {
		s.IsDefault = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("retrospective_process_matching_criteria"); ok {
		s.RetrospectiveProcessMatchingCriteria = value.(map[string]interface{})
	}

	res, err := c.CreateRetrospectiveProcess(s)
	if err != nil {
		return diag.Errorf("Error creating retrospective_process: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a retrospective_process resource: %s", d.Id()))

	return resourceRetrospectiveProcessRead(ctx, d, meta)
}

func resourceRetrospectiveProcessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading RetrospectiveProcess: %s", d.Id()))

	item, err := c.GetRetrospectiveProcess(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcess (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading retrospective_process: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("is_default", item.IsDefault)
	d.Set("retrospective_process_matching_criteria", item.RetrospectiveProcessMatchingCriteria)

	return nil
}

func resourceRetrospectiveProcessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating RetrospectiveProcess: %s", d.Id()))

	s := &client.RetrospectiveProcess{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("is_default") {
		s.IsDefault = tools.Bool(d.Get("is_default").(bool))
	}
	if d.HasChange("retrospective_process_matching_criteria") {
		s.RetrospectiveProcessMatchingCriteria = d.Get("retrospective_process_matching_criteria").(map[string]interface{})
	}

	_, err := c.UpdateRetrospectiveProcess(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating retrospective_process: %s", err.Error())
	}

	return resourceRetrospectiveProcessRead(ctx, d, meta)
}

func resourceRetrospectiveProcessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting RetrospectiveProcess: %s", d.Id()))

	err := c.DeleteRetrospectiveProcess(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("RetrospectiveProcess (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting retrospective_process: %s", err.Error())
	}

	d.SetId("")

	return nil
}
