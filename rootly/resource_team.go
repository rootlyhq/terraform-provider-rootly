package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Teams (e.g Infrastructure, Security, Search).",

		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the team",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the team",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"color": {
				Description:  "The color chosen for the team",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#E65252", // Default value from the API
				ValidateFunc: validCSSHexColor(),
			},
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating Team: %s", name))

	s := &client.Team{
		Name: name,
	}

	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}

	res, err := c.CreateTeam(s)
	if err != nil {
		return diag.Errorf("Error creating team: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a team resource: %v (%s)", name, d.Id()))

	return resourceTeamRead(ctx, d, meta)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Team: %s", d.Id()))

	team, err := c.GetTeam(d.Id())
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

	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("color", team.Color)

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Team: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.Team{
		Name: name,
	}

	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}

	if d.HasChange("color") {
		s.Color = d.Get("color").(string)
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
