package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
)

func resourceIncidentSubStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentSubStatusCreate,
		ReadContext:   resourceIncidentSubStatusRead,
		UpdateContext: resourceIncidentSubStatusUpdate,
		DeleteContext: resourceIncidentSubStatusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"incident_id": &schema.Schema{
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

			"assigned_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "",
			},

			"assigned_by_user_id": &schema.Schema{
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

func resourceIncidentSubStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentSubStatus"))

	s := &client.IncidentSubStatus{}

	if value, ok := d.GetOkExists("incident_id"); ok {
		s.IncidentId = value.(string)
	}
	if value, ok := d.GetOkExists("sub_status_id"); ok {
		s.SubStatusId = value.(string)
	}
	if value, ok := d.GetOkExists("assigned_at"); ok {
		s.AssignedAt = value.(string)
	}
	if value, ok := d.GetOkExists("assigned_by_user_id"); ok {
		s.AssignedByUserId = value.(int)
	}

	res, err := c.CreateIncidentSubStatus(s)
	if err != nil {
		return diag.Errorf("Error creating incident_sub_status: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_sub_status resource: %s", d.Id()))

	return resourceIncidentSubStatusRead(ctx, d, meta)
}

func resourceIncidentSubStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentSubStatus: %s", d.Id()))

	item, err := c.GetIncidentSubStatus(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentSubStatus (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_sub_status: %s", d.Id())
	}

	d.Set("incident_id", item.IncidentId)
	d.Set("sub_status_id", item.SubStatusId)
	d.Set("assigned_at", item.AssignedAt)
	d.Set("assigned_by_user_id", item.AssignedByUserId)

	return nil
}

func resourceIncidentSubStatusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentSubStatus: %s", d.Id()))

	s := &client.IncidentSubStatus{}

	if d.HasChange("incident_id") {
		s.IncidentId = d.Get("incident_id").(string)
	}
	if d.HasChange("sub_status_id") {
		s.SubStatusId = d.Get("sub_status_id").(string)
	}
	if d.HasChange("assigned_at") {
		s.AssignedAt = d.Get("assigned_at").(string)
	}
	if d.HasChange("assigned_by_user_id") {
		s.AssignedByUserId = d.Get("assigned_by_user_id").(int)
	}

	_, err := c.UpdateIncidentSubStatus(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_sub_status: %s", err.Error())
	}

	return resourceIncidentSubStatusRead(ctx, d, meta)
}

func resourceIncidentSubStatusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentSubStatus: %s", d.Id()))

	err := c.DeleteIncidentSubStatus(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentSubStatus (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_sub_status: %s", err.Error())
	}

	d.SetId("")

	return nil
}