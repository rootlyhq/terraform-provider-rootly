package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceSubStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubStatusCreate,
		ReadContext:   resourceSubStatusRead,
		UpdateContext: resourceSubStatusUpdate,
		DeleteContext: resourceSubStatusDelete,
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
				Description: "",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"parent_status": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "in_triage",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Value must be one of `in_triage`, `started`, `resolved`, `closed`, `cancelled`, `scheduled`, `in_progress`, `completed`.",
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

func resourceSubStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating SubStatus"))

	s := &client.SubStatus{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("parent_status"); ok {
		s.ParentStatus = value.(string)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}

	res, err := c.CreateSubStatus(s)
	if err != nil {
		return diag.Errorf("Error creating sub_status: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a sub_status resource: %s", d.Id()))

	return resourceSubStatusRead(ctx, d, meta)
}

func resourceSubStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading SubStatus: %s", d.Id()))

	item, err := c.GetSubStatus(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("SubStatus (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading sub_status: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("parent_status", item.ParentStatus)
	d.Set("position", item.Position)

	return nil
}

func resourceSubStatusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating SubStatus: %s", d.Id()))

	s := &client.SubStatus{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("parent_status") {
		s.ParentStatus = d.Get("parent_status").(string)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}

	_, err := c.UpdateSubStatus(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating sub_status: %s", err.Error())
	}

	return resourceSubStatusRead(ctx, d, meta)
}

func resourceSubStatusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting SubStatus: %s", d.Id()))

	err := c.DeleteSubStatus(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("SubStatus (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting sub_status: %s", err.Error())
	}

	d.SetId("")

	return nil
}