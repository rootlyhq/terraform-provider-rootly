package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceEscalationPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEscalationPathCreate,
		ReadContext:   resourceEscalationPathRead,
		UpdateContext: resourceEscalationPathUpdate,
		DeleteContext: resourceEscalationPathDelete,
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
				Description: "The name of the escalation path",
			},

			"default": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether this escalation path is the default path. Value must be one of true or false",
			},

			"notification_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Notification rule type",
			},

			"escalation_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the escalation policy",
			},

			"repeat": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Whether this path should be repeated until someone acknowledges the alert. Value must be one of true or false",
			},

			"repeat_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The number of times this path will be executed until someone acknowledges the alert",
			},

			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Escalation path rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"rule_type": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "alert_urgency",
							Required:    false,
							Optional:    true,
							ForceNew:    false,
							Description: "The type of the escalation path rule. Value must be one of `alert_urgency`, `working_hour`, `json_path`.",
						},

						"urgency_ids": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Computed:         true,
							Required:         false,
							Optional:         true,
							Description:      "Alert urgency ids for which this escalation path should be used",
						},

						"within_working_hour": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Description: "Whether the escalation path should be used within working hours. Value must be one of true or false",
						},

						"json_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							ForceNew:    false,
							Description: "JSON path to extract value from payload",
						},

						"operator": &schema.Schema{
							Type:        schema.TypeString,
							Default:     "is",
							Required:    false,
							Optional:    true,
							ForceNew:    false,
							Description: "How JSON path value should be matched. Value must be one of `is`, `is_not`, `contains`, `does_not_contain`.",
						},

						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							ForceNew:    false,
							Description: "Value with which JSON path value should be matched",
						},
					},
				},
			},
		},
	}
}

func resourceEscalationPathCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating EscalationPath"))

	s := &client.EscalationPath{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("default"); ok {
		s.Default = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("notification_type"); ok {
		s.NotificationType = value.(string)
	}
	if value, ok := d.GetOkExists("escalation_policy_id"); ok {
		s.EscalationPolicyId = value.(string)
	}
	if value, ok := d.GetOkExists("repeat"); ok {
		s.Repeat = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("repeat_count"); ok {
		s.RepeatCount = value.(int)
	}
	if value, ok := d.GetOkExists("rules"); ok {
		s.Rules = value.([]interface{})
	}

	res, err := c.CreateEscalationPath(s)
	if err != nil {
		return diag.Errorf("Error creating escalation_path: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a escalation_path resource: %s", d.Id()))

	return resourceEscalationPathRead(ctx, d, meta)
}

func resourceEscalationPathRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading EscalationPath: %s", d.Id()))

	item, err := c.GetEscalationPath(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationPath (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading escalation_path: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("default", item.Default)
	d.Set("notification_type", item.NotificationType)
	d.Set("escalation_policy_id", item.EscalationPolicyId)
	d.Set("repeat", item.Repeat)
	d.Set("repeat_count", item.RepeatCount)
	d.Set("rules", item.Rules)

	return nil
}

func resourceEscalationPathUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating EscalationPath: %s", d.Id()))

	s := &client.EscalationPath{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("default") {
		s.Default = tools.Bool(d.Get("default").(bool))
	}
	if d.HasChange("notification_type") {
		s.NotificationType = d.Get("notification_type").(string)
	}
	if d.HasChange("escalation_policy_id") {
		s.EscalationPolicyId = d.Get("escalation_policy_id").(string)
	}
	if d.HasChange("repeat") {
		s.Repeat = tools.Bool(d.Get("repeat").(bool))
	}
	if d.HasChange("repeat_count") {
		s.RepeatCount = d.Get("repeat_count").(int)
	}
	if d.HasChange("rules") {
		s.Rules = d.Get("rules").([]interface{})
	}

	_, err := c.UpdateEscalationPath(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating escalation_path: %s", err.Error())
	}

	return resourceEscalationPathRead(ctx, d, meta)
}

func resourceEscalationPathDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting EscalationPath: %s", d.Id()))

	err := c.DeleteEscalationPath(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EscalationPath (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting escalation_path: %s", err.Error())
	}

	d.SetId("")

	return nil
}