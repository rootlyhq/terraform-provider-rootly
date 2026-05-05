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

func resourceEdgeConnectorAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeConnectorActionCreate,
		ReadContext:   resourceEdgeConnectorActionRead,
		UpdateContext: resourceEdgeConnectorActionUpdate,
		DeleteContext: resourceEdgeConnectorActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"edge_connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the edge connector",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action name",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action slug",
			},
			"action_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "script",
				Description:  "Action type. Value must be one of `script`, `http`.",
				ValidateFunc: validation.StringInSlice([]string{"script", "http"}, false),
			},
			"icon": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Action icon. Value must be one of `bolt`, `bolt-slash`, `cog`, `command-line`, `code-bracket`, `server`, `server-stack`, `play`, `arrow-path`, `wrench-screwdriver`, `cube`, `rocket-launch`.",
				ValidateFunc: validation.StringInSlice([]string{"bolt", "bolt-slash", "cog", "command-line", "code-bracket", "server", "server-stack", "play", "arrow-path", "wrench-screwdriver", "cube", "rocket-launch"}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Action description",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout in seconds",
			},
			"parameters": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Description:      "Parameter definitions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "string",
							ValidateFunc: validation.StringInSlice([]string{"string", "number", "boolean", "list"}, false),
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"options": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Elem:             &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"last_executed_at": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceEdgeConnectorActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, "Creating EdgeConnectorAction")

	s := &client.EdgeConnectorAction{
		EdgeConnectorId: d.Get("edge_connector_id").(string),
	}

	if value, ok := d.GetOk("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOk("action_type"); ok {
		s.ActionType = value.(string)
	}
	if value, ok := d.GetOk("icon"); ok {
		s.Icon = value.(string)
	}
	if value, ok := d.GetOk("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOk("timeout"); ok {
		s.Timeout = value.(int)
	}
	if value, ok := d.GetOk("parameters"); ok {
		s.Metadata = map[string]interface{}{
			"parameters": flattenParameters(value.([]interface{})),
		}
	}

	res, err := c.CreateEdgeConnectorAction(s)
	if err != nil {
		return diag.Errorf("Error creating edge_connector_action: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a edge_connector_action resource: %s", d.Id()))

	return resourceEdgeConnectorActionRead(ctx, d, meta)
}

func resourceEdgeConnectorActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading EdgeConnectorAction: %s", d.Id()))

	item, err := c.GetEdgeConnectorAction(d.Get("edge_connector_id").(string), d.Id())
	if err != nil {
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EdgeConnectorAction (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading edge_connector_action: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("action_type", item.ActionType)
	d.Set("icon", item.Icon)
	d.Set("description", item.Description)
	d.Set("timeout", item.Timeout)
	d.Set("last_executed_at", item.LastExecutedAt)
	d.Set("created_at", item.CreatedAt)
	d.Set("updated_at", item.UpdatedAt)

	if item.Parameters != nil {
		d.Set("parameters", item.Parameters)
	}

	return nil
}

func resourceEdgeConnectorActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating EdgeConnectorAction: %s", d.Id()))

	s := &client.EdgeConnectorAction{
		EdgeConnectorId: d.Get("edge_connector_id").(string),
	}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("action_type") {
		s.ActionType = d.Get("action_type").(string)
	}
	if d.HasChange("icon") {
		s.Icon = d.Get("icon").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("timeout") {
		s.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("parameters") {
		s.Metadata = map[string]interface{}{
			"parameters": flattenParameters(d.Get("parameters").([]interface{})),
		}
	}

	_, err := c.UpdateEdgeConnectorAction(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating edge_connector_action: %s", err.Error())
	}

	return resourceEdgeConnectorActionRead(ctx, d, meta)
}

func resourceEdgeConnectorActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting EdgeConnectorAction: %s", d.Id()))

	err := c.DeleteEdgeConnectorAction(d.Get("edge_connector_id").(string), d.Id())
	if err != nil {
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("EdgeConnectorAction (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting edge_connector_action: %s", err.Error())
	}

	d.SetId("")
	return nil
}

// flattenParameters converts the Terraform list of parameter maps into the format the API expects.
func flattenParameters(params []interface{}) []interface{} {
	result := make([]interface{}, 0, len(params))
	for _, p := range params {
		if m, ok := p.(map[string]interface{}); ok {
			param := map[string]interface{}{}
			if v, ok := m["name"].(string); ok && v != "" {
				param["name"] = v
			}
			if v, ok := m["type"].(string); ok && v != "" {
				param["type"] = v
			}
			if v, ok := m["required"].(bool); ok {
				param["required"] = v
			}
			if v, ok := m["description"].(string); ok && v != "" {
				param["description"] = v
			}
			if v, ok := m["default"].(string); ok && v != "" {
				param["default"] = v
			}
			if opts, ok := m["options"].([]interface{}); ok && len(opts) > 0 {
				param["options"] = opts
			}
			result = append(result, param)
		}
	}
	return result
}
