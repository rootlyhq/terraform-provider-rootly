package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootly/terraform-provider-rootly/client"
)

func resourceCause() *schema.Resource {
	return &schema.Resource{
		Description: "Manages incident causes (e.g Bug, Load, Human Error, 3rd party Outage, Configuration Change).",

		CreateContext: resourceCauseCreate,
		ReadContext:   resourceCauseRead,
		UpdateContext: resourceCauseUpdate,
		DeleteContext: resourceCauseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the cause",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the cause",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCauseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Cause: %s", name))

	s := &client.Cause{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	res, err := c.CreateCause(s)
	if err != nil {
		return diag.Errorf("Error creating cause: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a cause resource: %v (%s)", name, d.Id()))

	return resourceCauseRead(ctx, d, meta)
}

func resourceCauseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Cause: %s", d.Id()))

	cause, err := c.GetCause(d.Id())
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

	d.Set("name", cause.Name)
	d.Set("description", cause.Description)

	return nil
}

func resourceCauseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Cause: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Cause{
		Name: name,
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
