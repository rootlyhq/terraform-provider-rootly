package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEscalationPolicyCreate,
		ReadContext:   resourceEscalationPolicyRead,
		UpdateContext: resourceEscalationPolicyUpdate,
		DeleteContext: resourceEscalationPolicyDelete,
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
				Description: "The name of the escalation policy",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the escalation policy",
			},

			"repeat_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The number of times this policy will be executed until someone acknowledges the alert",
			},

			"created_by_user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "User who created the escalation policy",
			},

			"last_updated_by_user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "User who updated the escalation policy",
			},

			"escalation_levels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Levels of the escalation policy",
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

			"services": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Services of the escalation policy",
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

			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Teams of to the escalation policy",
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

func resourceEscalationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating EscalationPolicy"))

	s := &client.EscalationPolicy{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("repeat_count"); ok {
		s.RepeatCount = value.(int)
	}
	if value, ok := d.GetOkExists("created_by_user_id"); ok {
		s.CreatedByUserId = value.(int)
	}
	if value, ok := d.GetOkExists("last_updated_by_user_id"); ok {
		s.LastUpdatedByUserId = value.(int)
	}
	if value, ok := d.GetOkExists("escalation_levels"); ok {
		s.EscalationLevels = value.([]interface{})
	}
	if value, ok := d.GetOkExists("services"); ok {
		s.Services = value.([]interface{})
	}
	if value, ok := d.GetOkExists("groups"); ok {
		s.Groups = value.([]interface{})
	}

	res, err := c.CreateEscalationPolicy(s)
	if err != nil {
		return diag.Errorf("Error creating escalation_policy: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a escalation_policy resource: %s", d.Id()))

	return resourceEscalationPolicyRead(ctx, d, meta)
}

func resourceEscalationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading EscalationPolicy: %s", d.Id()))

	item, err := c.GetEscalationPolicy(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationPolicy (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading escalation_policy: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("repeat_count", item.RepeatCount)
	d.Set("created_by_user_id", item.CreatedByUserId)
	d.Set("last_updated_by_user_id", item.LastUpdatedByUserId)
	d.Set("escalation_levels", item.EscalationLevels)
	d.Set("services", item.Services)
	d.Set("groups", item.Groups)

	return nil
}

func resourceEscalationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating EscalationPolicy: %s", d.Id()))

	s := &client.EscalationPolicy{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("repeat_count") {
		s.RepeatCount = d.Get("repeat_count").(int)
	}
	if d.HasChange("created_by_user_id") {
		s.CreatedByUserId = d.Get("created_by_user_id").(int)
	}
	if d.HasChange("last_updated_by_user_id") {
		s.LastUpdatedByUserId = d.Get("last_updated_by_user_id").(int)
	}
	if d.HasChange("escalation_levels") {
		s.EscalationLevels = d.Get("escalation_levels").([]interface{})
	}
	if d.HasChange("services") {
		s.Services = d.Get("services").([]interface{})
	}
	if d.HasChange("groups") {
		s.Groups = d.Get("groups").([]interface{})
	}

	_, err := c.UpdateEscalationPolicy(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating escalation_policy: %s", err.Error())
	}

	return resourceEscalationPolicyRead(ctx, d, meta)
}

func resourceEscalationPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting EscalationPolicy: %s", d.Id()))

	err := c.DeleteEscalationPolicy(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationPolicy (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting escalation_policy: %s", err.Error())
	}

	d.SetId("")

	return nil
}
