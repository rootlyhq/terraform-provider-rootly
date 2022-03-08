package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootly/terraform-provider-rootly/client"
)

func resourceSeverity() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceSeverityCreate,
		ReadContext:   resourceSeverityRead,
		UpdateContext: resourceSeverityUpdate,
		DeleteContext: resourceSeverityDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the severity",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the severity",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"severity": {
				Description: "The description of the severity",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"critical",
					"high",
					"medium",
					"low",
				}, false),
			},
			"color": {
				Description: "The color of the severity",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSeverityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)

	s := &client.Severity{
		Name: name,
	}

	if description, ok := d.GetOk("description"); ok {
		s.Description = description.(string)
	}

	if color, ok := d.GetOk("color"); ok {
		s.Color = color.(string)
	}

	severity, err := c.CreateSeverity(s)
	if err != nil {
		return diag.Errorf("Error creating severity: %s", err.Error())
	}

	d.SetId(severity.ID)
	tflog.Trace(ctx, "created a resource")

	return resourceSeverityRead(ctx, d, meta)
}

func resourceSeverityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	severity, err := c.GetSeverity(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Severity (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading severity: %s", d.Id())
	}

	d.Set("name", severity.Name)

	return diag.Errorf("not implemented")
}

func resourceSeverityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceSeverityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
