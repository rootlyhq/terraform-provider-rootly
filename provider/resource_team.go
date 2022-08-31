package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceTeam() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext: resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the team",
			},
			

			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The description of the team",
			},
			

				"notify_emails": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Emails to attach to the team",
				},
				

			"color": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			

				"slack_channels": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Slack Channels associated with this service",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
				},
				

				"slack_aliases": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Slack Aliases associated with this service",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
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
    if value, ok := d.GetOkExists("description"); ok {
				s.Description = value.(string)
			}
    if value, ok := d.GetOkExists("notify_emails"); ok {
				s.NotifyEmails = value.([]interface{})
			}
    if value, ok := d.GetOkExists("color"); ok {
				s.Color = value.(string)
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
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading team: %s", d.Id())
	}

	d.Set("name", item.Name)
  d.Set("description", item.Description)
  d.Set("notify_emails", item.NotifyEmails)
  d.Set("color", item.Color)
  d.Set("slack_channels", item.SlackChannels)
  d.Set("slack_aliases", item.SlackAliases)

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Team: %s", d.Id()))

	s := &client.Team{}

	  if d.HasChange("name") {
				s.Name = d.Get("name").(string)
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
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Team (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting team: %s", err.Error())
	}

	d.SetId("")

	return nil
}
