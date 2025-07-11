// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceWebhooksEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhooksEndpointCreate,
		ReadContext: resourceWebhooksEndpointRead,
		UpdateContext: resourceWebhooksEndpointUpdate,
		DeleteContext: resourceWebhooksEndpointDelete,
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
				Description: "The name of the endpoint",
				
			},
			

			"slug": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The slug of the endpoint",
				
			},
			

			"url": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The URL of the endpoint.",
				
			},
			

				"event_types": &schema.Schema {
					Type: schema.TypeList,
					Elem: &schema.Schema {
						Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"incident.created", "incident.updated", "incident.in_triage", "incident.mitigated", "incident.resolved", "incident.cancelled", "incident.deleted", "incident.scheduled.created", "incident.scheduled.updated", "incident.scheduled.in_progress", "incident.scheduled.completed", "incident.scheduled.deleted", "incident_post_mortem.created", "incident_post_mortem.updated", "incident_post_mortem.published", "incident_post_mortem.deleted", "incident_status_page_event.created", "incident_status_page_event.updated", "incident_status_page_event.deleted", "incident_event.created", "incident_event.updated", "incident_event.deleted", "alert.created", "pulse.created", "genius_workflow_run.queued", "genius_workflow_run.started", "genius_workflow_run.completed", "genius_workflow_run.failed", "genius_workflow_run.canceled"}, false),
					},
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Value must be one of `incident.created`, `incident.updated`, `incident.in_triage`, `incident.mitigated`, `incident.resolved`, `incident.cancelled`, `incident.deleted`, `incident.scheduled.created`, `incident.scheduled.updated`, `incident.scheduled.in_progress`, `incident.scheduled.completed`, `incident.scheduled.deleted`, `incident_post_mortem.created`, `incident_post_mortem.updated`, `incident_post_mortem.published`, `incident_post_mortem.deleted`, `incident_status_page_event.created`, `incident_status_page_event.updated`, `incident_status_page_event.deleted`, `incident_event.created`, `incident_event.updated`, `incident_event.deleted`, `alert.created`, `pulse.created`, `genius_workflow_run.queued`, `genius_workflow_run.started`, `genius_workflow_run.completed`, `genius_workflow_run.failed`, `genius_workflow_run.canceled`.",
					
				},
				

			"secret": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The webhook signing secret used to verify webhook requests.",
				
			},
			

				"enabled": &schema.Schema {
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
					
				},
				
		},
	}
}

func resourceWebhooksEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating WebhooksEndpoint"))

	s := &client.WebhooksEndpoint{}

	  if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("slug"); ok {
				s.Slug = value.(string)
			}
    if value, ok := d.GetOkExists("url"); ok {
				s.Url = value.(string)
			}
    if value, ok := d.GetOkExists("event_types"); ok {
				s.EventTypes = value.([]interface{})
			}
    if value, ok := d.GetOkExists("secret"); ok {
				s.Secret = value.(string)
			}
    if value, ok := d.GetOkExists("enabled"); ok {
				s.Enabled = tools.Bool(value.(bool))
			}

	res, err := c.CreateWebhooksEndpoint(s)
	if err != nil {
		return diag.Errorf("Error creating webhooks_endpoint: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a webhooks_endpoint resource: %s", d.Id()))

	return resourceWebhooksEndpointRead(ctx, d, meta)
}

func resourceWebhooksEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WebhooksEndpoint: %s", d.Id()))

	item, err := c.GetWebhooksEndpoint(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WebhooksEndpoint (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading webhooks_endpoint: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("slug", item.Slug)
  d.Set("url", item.Url)
  d.Set("event_types", item.EventTypes)
  d.Set("secret", item.Secret)
  d.Set("enabled", item.Enabled)

	return nil
}

func resourceWebhooksEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating WebhooksEndpoint: %s", d.Id()))

	s := &client.WebhooksEndpoint{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("slug") {
				s.Slug = d.Get("slug").(string)
			}
    if d.HasChange("url") {
				s.Url = d.Get("url").(string)
			}
    if d.HasChange("event_types") {
				s.EventTypes = d.Get("event_types").([]interface{})
			}
    if d.HasChange("secret") {
				s.Secret = d.Get("secret").(string)
			}
    if d.HasChange("enabled") {
				s.Enabled = tools.Bool(d.Get("enabled").(bool))
			}

	_, err := c.UpdateWebhooksEndpoint(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating webhooks_endpoint: %s", err.Error())
	}

	return resourceWebhooksEndpointRead(ctx, d, meta)
}

func resourceWebhooksEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting WebhooksEndpoint: %s", d.Id()))

	err := c.DeleteWebhooksEndpoint(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WebhooksEndpoint (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting webhooks_endpoint: %s", err.Error())
	}

	d.SetId("")

	return nil
}
