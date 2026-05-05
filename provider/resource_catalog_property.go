package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/rootlyhq/terraform-provider-rootly/v5/tools"
)

func resourceCatalogProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogPropertyCreate,
		ReadContext:   resourceCatalogPropertyRead,
		UpdateContext: resourceCatalogPropertyUpdate,
		DeleteContext: resourceCatalogPropertyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    true,
				WriteOnly:   false,
				Description: "",
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "",
			},

			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "text",
				Required:     false,
				Optional:     true,
				Sensitive:    false,
				ForceNew:     false,
				WriteOnly:    false,
				Description:  "Value must be one of `text`, `reference`.",
				ValidateFunc: validation.StringInSlice([]string{"text", "reference"}, false),
			},

			"kind_catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Restricts values to items of specified catalog.",
			},

			"multiple": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Whether the attribute accepts multiple values.. Value must be one of true or false",
			},

			"required": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Whether the property is required.. Value must be one of true or false",
			},

			"catalog_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "catalog",
				Required:     false,
				Optional:     true,
				Sensitive:    false,
				ForceNew:     false,
				WriteOnly:    false,
				Description:  "The type of catalog the property belongs to.. Value must be one of `catalog`, `cause`, `environment`, `functionality`, `incident_type`, `service`, `team`.",
				ValidateFunc: validation.StringInSlice([]string{"catalog", "cause", "environment", "functionality", "incident_type", "service", "team"}, false),
			},
		},
	}
}

func resourceCatalogPropertyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating CatalogProperty"))

	s := &client.CatalogProperty{}

	if value, ok := d.GetOkExists("catalog_id"); ok {
		s.CatalogId = value.(string)
	}
	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("kind"); ok {
		s.Kind = value.(string)
	}
	if value, ok := d.GetOkExists("kind_catalog_id"); ok {
		s.KindCatalogId = value.(string)
	}
	if value, ok := d.GetOkExists("multiple"); ok {
		s.Multiple = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("required"); ok {
		s.Required = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("catalog_type"); ok {
		s.CatalogType = value.(string)
	}

	res, err := c.CreateCatalogProperty(s)
	if err != nil {
		return diag.Errorf("Error creating catalog_property: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a catalog_property resource: %s", d.Id()))

	return resourceCatalogPropertyRead(ctx, d, meta)
}

func resourceCatalogPropertyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading CatalogProperty: %s", d.Id()))

	item, err := c.GetCatalogProperty(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CatalogProperty (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading catalog_property: %s", d.Id())
	}

	d.Set("catalog_id", item.CatalogId)
	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("kind", item.Kind)
	d.Set("kind_catalog_id", item.KindCatalogId)
	d.Set("multiple", item.Multiple)
	d.Set("required", item.Required)
	d.Set("catalog_type", item.CatalogType)

	return nil
}

func resourceCatalogPropertyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating CatalogProperty: %s", d.Id()))

	s := &client.CatalogProperty{}

	if d.HasChange("catalog_id") {
		s.CatalogId = d.Get("catalog_id").(string)
	}
	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("kind") {
		s.Kind = d.Get("kind").(string)
	}
	if d.HasChange("kind_catalog_id") {
		s.KindCatalogId = d.Get("kind_catalog_id").(string)
	}
	if d.HasChange("multiple") {
		s.Multiple = tools.Bool(d.Get("multiple").(bool))
	}
	if d.HasChange("required") {
		s.Required = tools.Bool(d.Get("required").(bool))
	}
	if d.HasChange("catalog_type") {
		s.CatalogType = d.Get("catalog_type").(string)
	}

	_, err := c.UpdateCatalogProperty(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating catalog_property: %s", err.Error())
	}

	return resourceCatalogPropertyRead(ctx, d, meta)
}

func resourceCatalogPropertyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting CatalogProperty: %s", d.Id()))

	err := c.DeleteCatalogProperty(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CatalogProperty (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting catalog_property: %s", err.Error())
	}

	d.SetId("")

	return nil
}
