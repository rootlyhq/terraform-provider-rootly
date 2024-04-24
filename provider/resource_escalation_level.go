package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceEscalationLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEscalationLevelCreate,
		ReadContext:   resourceEscalationLevelRead,
		UpdateContext: resourceEscalationLevelUpdate,
		DeleteContext: resourceEscalationLevelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"escalation_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the escalation policy",
			},

			"delay": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Delay before notification targets will be alerted.",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Position of the escalation policy level",
			},

			"notification_target_params": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Escalation level's notification targets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceEscalationLevelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating EscalationLevel"))

	s := &client.EscalationLevel{}

	if value, ok := d.GetOkExists("escalation_policy_id"); ok {
		s.EscalationPolicyId = value.(string)
	}
	if value, ok := d.GetOkExists("delay"); ok {
		s.Delay = value.(int)
	}
	if value, ok := d.GetOkExists("position"); ok {
		s.Position = value.(int)
	}
	if value, ok := d.GetOkExists("notification_target_params"); ok {
		s.NotificationTargetParams = value.([]interface{})
	}

	res, err := c.CreateEscalationLevel(s)
	if err != nil {
		return diag.Errorf("Error creating escalation_level: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a escalation_level resource: %s", d.Id()))

	return resourceEscalationLevelRead(ctx, d, meta)
}

func resourceEscalationLevelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading EscalationLevel: %s", d.Id()))

	item, err := c.GetEscalationLevel(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationLevel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading escalation_level: %s", d.Id())
	}

	d.Set("escalation_policy_id", item.EscalationPolicyId)
	d.Set("delay", item.Delay)
	d.Set("position", item.Position)
	d.Set("notification_target_params", item.NotificationTargetParams)

	return nil
}

func resourceEscalationLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating EscalationLevel: %s", d.Id()))

	s := &client.EscalationLevel{}

	if d.HasChange("escalation_policy_id") {
		s.EscalationPolicyId = d.Get("escalation_policy_id").(string)
	}
	if d.HasChange("delay") {
		s.Delay = d.Get("delay").(int)
	}
	if d.HasChange("position") {
		s.Position = d.Get("position").(int)
	}
	if d.HasChange("notification_target_params") {
		s.NotificationTargetParams = d.Get("notification_target_params").([]interface{})
	}

	_, err := c.UpdateEscalationLevel(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating escalation_level: %s", err.Error())
	}

	return resourceEscalationLevelRead(ctx, d, meta)
}

func resourceEscalationLevelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting EscalationLevel: %s", d.Id()))

	err := c.DeleteEscalationLevel(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationLevel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting escalation_level: %s", err.Error())
	}

	d.SetId("")

	return nil
}
