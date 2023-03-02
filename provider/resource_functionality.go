package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceFunctionality() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionalityCreate,
		ReadContext:   resourceFunctionalityRead,
		UpdateContext: resourceFunctionalityUpdate,
		DeleteContext: resourceFunctionalityDelete,
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
				Description: "The name of the functionality",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The slug of the functionality",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the functionality",
			},

			"public_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The public description of the functionality",
			},

			"notify_emails": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Emails attached to the functionality",
			},

			"color": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "",
			},

			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Position of the functionality",
			},

			"environment_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Environments associated with this functionality",
			},

			"service_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Services associated with this functionality",
			},

			"owners_group_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Owner Teams associated with this functionality",
			},

			"owners_user_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Owner Users associated with this service",
			},

			"slack_channels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Slack Channels associated with this service",
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

			"slack_aliases": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Slack Aliases associated with this service",
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

func resourceFunctionalityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Functionality"))

	s := &client.Functionality{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("public_description"); ok {
		s.PublicDescription = value.(string)
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
	if value, ok := d.GetOkExists("environment_ids"); ok {
		s.EnvironmentIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("service_ids"); ok {
		s.ServiceIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("owners_group_ids"); ok {
		s.OwnersGroupIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("owners_user_ids"); ok {
		s.OwnersUserIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("slack_channels"); ok {
		s.SlackChannels = value.([]interface{})
	}
	if value, ok := d.GetOkExists("slack_aliases"); ok {
		s.SlackAliases = value.([]interface{})
	}

	res, err := c.CreateFunctionality(s)
	if err != nil {
		return diag.Errorf("Error creating functionality: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a functionality resource: %s", d.Id()))

	return resourceFunctionalityRead(ctx, d, meta)
}

func resourceFunctionalityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Functionality: %s", d.Id()))

	item, err := c.GetFunctionality(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Functionality (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading functionality: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("public_description", item.PublicDescription)
	d.Set("notify_emails", item.NotifyEmails)
	d.Set("color", item.Color)
	d.Set("position", item.Position)
	d.Set("environment_ids", item.EnvironmentIds)
	d.Set("service_ids", item.ServiceIds)
	d.Set("owners_group_ids", item.OwnersGroupIds)
	d.Set("owners_user_ids", item.OwnersUserIds)
	d.Set("slack_channels", item.SlackChannels)
	d.Set("slack_aliases", item.SlackAliases)

	return nil
}

func resourceFunctionalityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Functionality: %s", d.Id()))

	s := &client.Functionality{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("public_description") {
		s.PublicDescription = d.Get("public_description").(string)
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
	if d.HasChange("environment_ids") {
		s.EnvironmentIds = d.Get("environment_ids").([]interface{})
	}
	if d.HasChange("service_ids") {
		s.ServiceIds = d.Get("service_ids").([]interface{})
	}
	if d.HasChange("owners_group_ids") {
		s.OwnersGroupIds = d.Get("owners_group_ids").([]interface{})
	}
	if d.HasChange("owners_user_ids") {
		s.OwnersUserIds = d.Get("owners_user_ids").([]interface{})
	}
	if d.HasChange("slack_channels") {
		s.SlackChannels = d.Get("slack_channels").([]interface{})
	}
	if d.HasChange("slack_aliases") {
		s.SlackAliases = d.Get("slack_aliases").([]interface{})
	}

	_, err := c.UpdateFunctionality(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating functionality: %s", err.Error())
	}

	return resourceFunctionalityRead(ctx, d, meta)
}

func resourceFunctionalityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Functionality: %s", d.Id()))

	err := c.DeleteFunctionality(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Functionality (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting functionality: %s", err.Error())
	}

	d.SetId("")

	return nil
}
