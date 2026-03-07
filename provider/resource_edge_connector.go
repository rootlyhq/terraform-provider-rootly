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

func resourceEdgeConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeConnectorCreate,
		ReadContext:   resourceEdgeConnectorRead,
		UpdateContext: resourceEdgeConnectorUpdate,
		DeleteContext: resourceEdgeConnectorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connector name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Connector description",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Connector status. Value must be one of `active`, `paused`.",
				ValidateFunc: validation.StringInSlice([]string{"active", "paused"}, false),
			},
			"subscriptions": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description:      "Array of event types to subscribe to",
			},
			"online": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether connector is currently online",
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEdgeConnectorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, "Creating EdgeConnector")

	s := &client.EdgeConnector{}

	if value, ok := d.GetOk("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOk("status"); ok {
		s.Status = value.(string)
	}

	res, err := c.CreateEdgeConnector(s)
	if err != nil {
		return diag.Errorf("Error creating edge_connector: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a edge_connector resource: %s", d.Id()))

	return resourceEdgeConnectorRead(ctx, d, meta)
}

func resourceEdgeConnectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading EdgeConnector: %s", d.Id()))

	item, err := c.GetEdgeConnector(d.Id())
	if err != nil {
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EdgeConnector (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading edge_connector: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("status", item.Status)
	d.Set("created_at", item.CreatedAt)
	d.Set("updated_at", item.UpdatedAt)

	return nil
}

func resourceEdgeConnectorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating EdgeConnector: %s", d.Id()))

	s := &client.EdgeConnector{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("status") {
		s.Status = d.Get("status").(string)
	}

	_, err := c.UpdateEdgeConnector(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating edge_connector: %s", err.Error())
	}

	return resourceEdgeConnectorRead(ctx, d, meta)
}

func resourceEdgeConnectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting EdgeConnector: %s", d.Id()))

	err := c.DeleteEdgeConnector(d.Id())
	if err != nil {
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EdgeConnector (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting edge_connector: %s", err.Error())
	}

	d.SetId("")
	return nil
}
