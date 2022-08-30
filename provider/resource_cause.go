package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceCause() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceCauseCreate,
		ReadContext: resourceCauseRead,
		UpdateContext: resourceCauseUpdate,
		DeleteContext: resourceCauseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The name of the cause",
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The slug of the cause",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "The description of the cause",
			},
			
		},
	}
}

func resourceCauseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Cause"))

	s := &client.Cause{}

	  if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
    if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
    if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}

	res, err := c.CreateCause(s)
	if err != nil {
		return diag.Errorf("Error creating cause: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a cause resource: %s", d.Id()))

	return resourceCauseRead(ctx, d, meta)
}

func resourceCauseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Cause: %s", d.Id()))

	item, err := c.GetCause(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Cause (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading cause: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)

	return nil
}

func resourceCauseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Cause: %s", d.Id()))

	s := &client.Cause{}

	  if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
    if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
    if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	_, err := c.UpdateCause(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating cause: %s", err.Error())
	}

	return resourceCauseRead(ctx, d, meta)
}

func resourceCauseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Cause: %s", d.Id()))

	err := c.DeleteCause(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Cause (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting cause: %s", err.Error())
	}

	d.SetId("")

	return nil
}
