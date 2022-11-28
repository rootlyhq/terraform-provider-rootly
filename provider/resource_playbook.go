package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourcePlaybook() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourcePlaybookCreate,
		ReadContext: resourcePlaybookRead,
		UpdateContext: resourcePlaybookUpdate,
		DeleteContext: resourcePlaybookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"title": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The title of the playbook",
			},
			

			"summary": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The summary of the playbook",
			},
			

			"external_url": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The external url of the playbook",
			},
			
		},
	}
}


func resourcePlaybookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Playbook"))

	s := &client.Playbook{}

	  if value, ok := d.GetOkExists("title"); ok {
				s.Title = value.(string)
			}
    if value, ok := d.GetOkExists("summary"); ok {
				s.Summary = value.(string)
			}
    if value, ok := d.GetOkExists("external_url"); ok {
				s.ExternalUrl = value.(string)
			}

	res, err := c.CreatePlaybook(s)
	if err != nil {
		return diag.Errorf("Error creating playbook: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a playbook resource: %s", d.Id()))

	return resourcePlaybookRead(ctx, d, meta)
}


func resourcePlaybookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Playbook: %s", d.Id()))

	item, err := c.GetPlaybook(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Playbook (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading playbook: %s", d.Id())
	}

	d.Set("title", item.Title)
  d.Set("summary", item.Summary)
  d.Set("external_url", item.ExternalUrl)

	return nil
}


func resourcePlaybookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Playbook: %s", d.Id()))

	s := &client.Playbook{}

	  if d.HasChange("title") {
				s.Title = d.Get("title").(string)
			}
    if d.HasChange("summary") {
				s.Summary = d.Get("summary").(string)
			}
    if d.HasChange("external_url") {
				s.ExternalUrl = d.Get("external_url").(string)
			}

	_, err := c.UpdatePlaybook(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating playbook: %s", err.Error())
	}

	return resourcePlaybookRead(ctx, d, meta)
}


func resourcePlaybookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Playbook: %s", d.Id()))

	err := c.DeletePlaybook(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Playbook (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting playbook: %s", err.Error())
	}

	d.SetId("")

	return nil
}

