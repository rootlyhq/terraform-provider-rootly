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

func resourceIncidentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentTypeCreate,
		ReadContext: resourceIncidentTypeRead,
		UpdateContext: resourceIncidentTypeUpdate,
		DeleteContext: resourceIncidentTypeDelete,
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
				Description: "The name of the incident type",
				
			},
			

			"slug": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the incident type",
				
			},
			

			"description": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the incident type",
				
			},
			

			"color": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The hex color of the incident type",
				
			},
			

		"position": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "Position of the incident type",
			
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
					Description: "Emails to attach to the incident type",
					
				},
				

				"slack_channels": &schema.Schema {
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Slack Channels associated with this incident type",
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
					Description: "Slack Aliases associated with this incident type",
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

func resourceIncidentTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentType"))

	s := &client.IncidentType{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("color"); ok {
				s.Color = value.(string)
			}
    if value, ok := d.GetOkExists("position"); ok {
				s.Position = value.(int)
			}
    if value, ok := d.GetOkExists("notify_emails"); ok {
				s.NotifyEmails = value.([]interface{})
			}
    if value, ok := d.GetOkExists("slack_channels"); ok {
				s.SlackChannels = value.([]interface{})
			}
    if value, ok := d.GetOkExists("slack_aliases"); ok {
				s.SlackAliases = value.([]interface{})
			}

	res, err := c.CreateIncidentType(s)
	if err != nil {
		return diag.Errorf("Error creating incident_type: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_type resource: %s", d.Id()))

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentType: %s", d.Id()))

	item, err := c.GetIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_type: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("description", item.Description)
  d.Set("color", item.Color)
  d.Set("position", item.Position)
  d.Set("notify_emails", item.NotifyEmails)
  
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

func resourceIncidentTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentType: %s", d.Id()))

	s := &client.IncidentType{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("description") {
				s.Description = d.Get("description").(string)
			}
    if d.HasChange("color") {
				s.Color = d.Get("color").(string)
			}
    if d.HasChange("position") {
				s.Position = d.Get("position").(int)
			}
    if d.HasChange("notify_emails") {
				s.NotifyEmails = d.Get("notify_emails").([]interface{})
			}
    if d.HasChange("slack_channels") {
				s.SlackChannels = d.Get("slack_channels").([]interface{})
			}
    if d.HasChange("slack_aliases") {
				s.SlackAliases = d.Get("slack_aliases").([]interface{})
			}

	_, err := c.UpdateIncidentType(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_type: %s", err.Error())
	}

	return resourceIncidentTypeRead(ctx, d, meta)
}

func resourceIncidentTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentType: %s", d.Id()))

	err := c.DeleteIncidentType(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentType (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_type: %s", err.Error())
	}

	d.SetId("")

	return nil
}
