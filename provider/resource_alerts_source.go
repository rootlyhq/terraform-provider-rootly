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

func resourceAlertsSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertsSourceCreate,
		ReadContext:   resourceAlertsSourceRead,
		UpdateContext: resourceAlertsSourceUpdate,
		DeleteContext: resourceAlertsSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"alert_urgency_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The alert urgency ID",
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the alert source",
			},

			"source_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The alert source type",
			},

			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The current status of the alert source",
			},

			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "A secret key used to authenticate incoming requests to this alerts source",
			},

			"webhook_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The URL endpoint of the alert source",
			},

			"alert_source_urgency_rules_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_urgency_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
						"json_path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
						"operator": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
					},
				},
			},
			"alert_template_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
						"external_url": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
						},
					},
				},
			},

			"sourceable_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_resolve": &schema.Schema{
							Type:     schema.TypeBool,
							Default:  false,
							Required: false,
							Optional: true,
						},

						"resolve_state": &schema.Schema{
							Type:        schema.TypeString,
							Required:    false,
							Optional:    true,
							Description: "This value is matched with the value extracted from alerts payload using JSON path in field_mappings_attributes",
						},

						"field_mappings_attributes": &schema.Schema{
							Type:             schema.TypeList,
							Optional:         true,
							MinItems:         0,
							MaxItems:         25,
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": &schema.Schema{
										Type:     schema.TypeString,
										Required: false,
										Optional: true,
									},
									"json_path": &schema.Schema{
										Type:     schema.TypeString,
										Required: false,
										Optional: true,
									},
								},
							},
						},
					},
				},
				Description: "Additional attributes specific to certain alert sources (e.g., generic_webhook), encapsulating source-specific configurations or details",
			},
		},
	}
}

func resourceAlertsSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating AlertsSource"))

	s := &client.AlertsSource{}

	if value, ok := d.GetOkExists("alert_urgency_id"); ok {
		s.AlertUrgencyId = value.(string)
	}
	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("source_type"); ok {
		s.SourceType = value.(string)
	}
	if value, ok := d.GetOkExists("status"); ok {
		s.Status = value.(string)
	}
	if value, ok := d.GetOkExists("secret"); ok {
		s.Secret = value.(string)
	}
	if value, ok := d.GetOkExists("webhook_endpoint"); ok {
		s.WebhookEndpoint = value.(string)
	}
	if value, ok := d.GetOkExists("alert_source_urgency_rules_attributes"); ok {
		s.AlertSourceUrgencyRulesAttributes = value.([]interface{})
	}
	if value, ok := d.GetOkExists("sourceable_attributes"); ok {
		s.SourceableAttributes = value.([]interface{})[0].(map[string]interface{})
	}
	if value, ok := d.GetOkExists("alert_template_attributes"); ok {
		s.AlertTemplateAttributes = value.([]interface{})[0].(map[string]interface{})
	}

	res, err := c.CreateAlertsSource(s)
	if err != nil {
		return diag.Errorf("Error creating alerts_source: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a alerts_source resource: %s", d.Id()))

	return resourceAlertsSourceRead(ctx, d, meta)
}

func resourceAlertsSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading AlertsSource: %s", d.Id()))

	item, err := c.GetAlertsSource(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertsSource (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading alerts_source: %s", d.Id())
	}

	d.Set("alert_urgency_id", item.AlertUrgencyId)
	d.Set("name", item.Name)
	d.Set("source_type", item.SourceType)
	d.Set("status", item.Status)
	d.Set("secret", item.Secret)
	d.Set("webhook_endpoint", item.WebhookEndpoint)
	d.Set("alert_source_urgency_rules_attributes", item.AlertSourceUrgencyRulesAttributes)
	if item_field_mappings, ok := item.SourceableAttributes["field_mappings_attributes"].([]interface{}); ok {
		sourceables := make([]interface{}, 1, 1)
		field_mappings := make([]map[string]interface{}, 0)
		for _, fma := range item_field_mappings {
			field_mapping := make(map[string]interface{}, 0)
			if fmaMap, ok := fma.(map[string]interface{}); ok {
				for key, val := range fmaMap {
					if key == "field" || key == "json_path" {
						field_mapping[key] = val
					}
				}
			}
			field_mappings = append(field_mappings, field_mapping)
		}
		sourceables[0] = map[string]interface{}{
			"auto_resolve":              item.SourceableAttributes["auto_resolve"],
			"resolve_state":             item.SourceableAttributes["resolve_state"],
			"field_mappings_attributes": field_mappings,
		}
		d.Set("sourceable_attributes", sourceables)
	}

	if item.AlertTemplateAttributes["title"] != nil || item.AlertTemplateAttributes["description"] != nil || item.AlertTemplateAttributes["external_url"] != nil {
		alert_templates := make([]interface{}, 1, 1)
		alert_templates[0] = map[string]interface{}{
			"title":        item.AlertTemplateAttributes["title"],
			"description":  item.AlertTemplateAttributes["description"],
			"external_url": item.AlertTemplateAttributes["external_url"],
		}
		d.Set("alert_template_attributes", alert_templates)
	}

	return nil
}

func resourceAlertsSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating AlertsSource: %s", d.Id()))

	s := &client.AlertsSource{}

	if d.HasChange("alert_urgency_id") {
		s.AlertUrgencyId = d.Get("alert_urgency_id").(string)
	}
	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("source_type") {
		s.SourceType = d.Get("source_type").(string)
	}
	if d.HasChange("status") {
		s.Status = d.Get("status").(string)
	}
	if d.HasChange("secret") {
		s.Secret = d.Get("secret").(string)
	}
	if d.HasChange("webhook_endpoint") {
		s.WebhookEndpoint = d.Get("webhook_endpoint").(string)
	}
	if d.HasChange("alert_source_urgency_rules_attributes") {
		s.AlertSourceUrgencyRulesAttributes = d.Get("alert_source_urgency_rules_attributes").([]interface{})
	}
	if d.HasChange("alert_template_attributes") {
		tps := d.Get("alert_template_attributes").([]interface{})
		for _, tpsi := range tps {
			s.AlertTemplateAttributes = tpsi.(map[string]interface{})
		}
	}
	if d.HasChange("sourceable_attributes") {
		tps := d.Get("sourceable_attributes").([]interface{})
		for _, tpsi := range tps {
			s.SourceableAttributes = tpsi.(map[string]interface{})
		}
	}

	_, err := c.UpdateAlertsSource(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating alerts_source: %s", err.Error())
	}

	return resourceAlertsSourceRead(ctx, d, meta)
}

func resourceAlertsSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting AlertsSource: %s", d.Id()))

	err := c.DeleteAlertsSource(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("AlertsSource (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting alerts_source: %s", err.Error())
	}

	d.SetId("")

	return nil
}
