// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext: resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter {
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema {
			
			"name": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the team",
				
			},
			

			"slug": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
				
			},
			

			"description": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the team",
				
			},
			

				"notify_emails": &schema.Schema {
					Type: schema.TypeList,
					Elem: &schema.Schema {
						Type: schema.TypeString,
					},
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Emails to attach to the team",
					
				},
				

			"color": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The hex color of the team",
				
			},
			

		"position": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "Position of the team",
			
		},
		

			"backstage_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The Backstage entity id associated to this team. eg: :namespace/:kind/:entity_name",
				
			},
			

			"external_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The external id associated to this team",
				
			},
			

			"pagerduty_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The PagerDuty group id associated to this team",
				
			},
			

			"pagerduty_service_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The PagerDuty service id associated to this team",
				
			},
			

			"opsgenie_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The Opsgenie group id associated to this team",
				
			},
			

			"victor_ops_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The VictorOps group id associated to this team",
				
			},
			

			"pagertree_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The PagerTree group id associated to this team",
				
			},
			

			"cortex_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The Cortex group id associated to this team",
				
			},
			

			"service_now_ci_sys_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The Service Now CI sys id associated to this team",
				
			},
			

				"user_ids": &schema.Schema {
					Type: schema.TypeList,
					Elem: &schema.Schema {
						Type: schema.TypeInt,
					},
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "The user ids of the members of this team.",
					
				},
				

				"admin_ids": &schema.Schema {
					Type: schema.TypeList,
					Elem: &schema.Schema {
						Type: schema.TypeInt,
					},
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "The user ids of the admins of this team. These users must also be present in user_ids attribute.",
					
				},
				

			"alerts_email_enabled": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "Enable alerts through email. Value must be one of true or false",
				
			},
			

			"alerts_email_address": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "Email generated to send alerts to",
				
			},
			

			"alert_urgency_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The alert urgency id of the team",
				
			},
			

				"slack_channels": &schema.Schema {
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Slack Channels associated with this team",
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Elem: &schema.Resource {
						Schema: map[string]*schema.Schema {
              
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "Slack channel ID",
				
			},
			

			"name": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "Slack channel name",
				
			},
			
						},
					},
					
				},
				

				"slack_aliases": &schema.Schema {
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Slack Aliases associated with this team",
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Elem: &schema.Resource {
						Schema: map[string]*schema.Schema {
              
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "Slack alias ID",
				
			},
			

			"name": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "Slack alias name",
				
			},
			
						},
					},
					
				},
				
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Team"))

	s := &client.Team{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("notify_emails"); ok {
				s.NotifyEmails = value.([]interface{})
			}
    if value, ok := d.GetOkExists("color"); ok {
				s.Color = value.(string)
			}
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}
    if value, ok := d.GetOkExists("backstage_id"); ok {
				s.BackstageId = value.(string)
			}
    if value, ok := d.GetOkExists("external_id"); ok {
				s.ExternalId = value.(string)
			}
    if value, ok := d.GetOkExists("pagerduty_id"); ok {
				s.PagerdutyId = value.(string)
			}
    if value, ok := d.GetOkExists("pagerduty_service_id"); ok {
				s.PagerdutyServiceId = value.(string)
			}
    if value, ok := d.GetOkExists("opsgenie_id"); ok {
				s.OpsgenieId = value.(string)
			}
    if value, ok := d.GetOkExists("victor_ops_id"); ok {
				s.VictorOpsId = value.(string)
			}
    if value, ok := d.GetOkExists("pagertree_id"); ok {
				s.PagertreeId = value.(string)
			}
    if value, ok := d.GetOkExists("cortex_id"); ok {
				s.CortexId = value.(string)
			}
    if value, ok := d.GetOkExists("service_now_ci_sys_id"); ok {
				s.ServiceNowCiSysId = value.(string)
			}
    if value, ok := d.GetOkExists("user_ids"); ok {
				s.UserIds = value.([]interface{})
			}
    if value, ok := d.GetOkExists("admin_ids"); ok {
				s.AdminIds = value.([]interface{})
			}
    if value, ok := d.GetOkExists("alerts_email_enabled"); ok {
				s.AlertsEmailEnabled = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("alerts_email_address"); ok {
				s.AlertsEmailAddress = value.(string)
			}
    if value, ok := d.GetOkExists("alert_urgency_id"); ok {
				s.AlertUrgencyId = value.(string)
			}
    if value, ok := d.GetOkExists("slack_channels"); ok {
				s.SlackChannels = value.([]interface{})
			}
    if value, ok := d.GetOkExists("slack_aliases"); ok {
				s.SlackAliases = value.([]interface{})
			}

	res, err := c.CreateTeam(s)
	if err != nil {
		return diag.Errorf("Error creating team: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a team resource: %s", d.Id()))

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Team: %s", d.Id()))

	item, err := c.GetTeam(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading team: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("notify_emails", item.NotifyEmails)
  d.Set("color", item.Color)
  d.Set("position", item.Position)
  d.Set("backstage_id", item.BackstageId)
  d.Set("external_id", item.ExternalId)
  d.Set("pagerduty_id", item.PagerdutyId)
  d.Set("pagerduty_service_id", item.PagerdutyServiceId)
  d.Set("opsgenie_id", item.OpsgenieId)
  d.Set("victor_ops_id", item.VictorOpsId)
  d.Set("pagertree_id", item.PagertreeId)
  d.Set("cortex_id", item.CortexId)
  d.Set("service_now_ci_sys_id", item.ServiceNowCiSysId)
  d.Set("user_ids", item.UserIds)
  d.Set("admin_ids", item.AdminIds)
  d.Set("alerts_email_enabled", item.AlertsEmailEnabled)
  d.Set("alerts_email_address", item.AlertsEmailAddress)
  d.Set("alert_urgency_id", item.AlertUrgencyId)
  
          if item.SlackChannels != nil {
              processedItems := make([]map[string]interface{}, 0)

              for _, c := range item.SlackChannels {
                  if rawItem, ok := c.(map[string]interface{}); ok {
                      // Create a new map with only the fields defined in the schema
                      processedItem := map[string]interface{}{
                          "id": rawItem["id"],
"name": rawItem["name"],
                      }
                      processedItems = append(processedItems, processedItem)
                  }
              }

              d.Set("slack_channels", processedItems)
          } else {
              d.Set("slack_channels", nil)
          }
        
  
          if item.SlackAliases != nil {
              processedItems := make([]map[string]interface{}, 0)

              for _, c := range item.SlackAliases {
                  if rawItem, ok := c.(map[string]interface{}); ok {
                      // Create a new map with only the fields defined in the schema
                      processedItem := map[string]interface{}{
                          "id": rawItem["id"],
"name": rawItem["name"],
                      }
                      processedItems = append(processedItems, processedItem)
                  }
              }

              d.Set("slack_aliases", processedItems)
          } else {
              d.Set("slack_aliases", nil)
          }
        

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Team: %s", d.Id()))

	s := &client.Team{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}
    if d.HasChange("notify_emails") {
				s.NotifyEmails = d.Get("notify_emails").([]interface{})
			}
    if d.HasChange("color") {
				s.Color = d.Get("color").(string)
			}
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}
    if d.HasChange("backstage_id") {
				s.BackstageId = d.Get("backstage_id").(string)
			}
    if d.HasChange("external_id") {
				s.ExternalId = d.Get("external_id").(string)
			}
    if d.HasChange("pagerduty_id") {
				s.PagerdutyId = d.Get("pagerduty_id").(string)
			}
    if d.HasChange("pagerduty_service_id") {
				s.PagerdutyServiceId = d.Get("pagerduty_service_id").(string)
			}
    if d.HasChange("opsgenie_id") {
				s.OpsgenieId = d.Get("opsgenie_id").(string)
			}
    if d.HasChange("victor_ops_id") {
				s.VictorOpsId = d.Get("victor_ops_id").(string)
			}
    if d.HasChange("pagertree_id") {
				s.PagertreeId = d.Get("pagertree_id").(string)
			}
    if d.HasChange("cortex_id") {
				s.CortexId = d.Get("cortex_id").(string)
			}
    if d.HasChange("service_now_ci_sys_id") {
				s.ServiceNowCiSysId = d.Get("service_now_ci_sys_id").(string)
			}
    if d.HasChange("user_ids") {
				s.UserIds = d.Get("user_ids").([]interface{})
			}
    if d.HasChange("admin_ids") {
				s.AdminIds = d.Get("admin_ids").([]interface{})
			}
    if d.HasChange("alerts_email_enabled") {
				s.AlertsEmailEnabled = tools.Bool(d.Get("alerts_email_enabled").(bool))
			}
    if d.HasChange("alerts_email_address") {
				s.AlertsEmailAddress = d.Get("alerts_email_address").(string)
			}
    if d.HasChange("alert_urgency_id") {
				s.AlertUrgencyId = d.Get("alert_urgency_id").(string)
			}
    if d.HasChange("slack_channels") {
				s.SlackChannels = d.Get("slack_channels").([]interface{})
			}
    if d.HasChange("slack_aliases") {
				s.SlackAliases = d.Get("slack_aliases").([]interface{})
			}

	_, err := c.UpdateTeam(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating team: %s", err.Error())
	}

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Team: %s", d.Id()))

	err := c.DeleteTeam(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting team: %s", err.Error())
	}

	d.SetId("")

	return nil
}
