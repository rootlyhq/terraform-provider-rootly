package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func resourceDashboardPanel() *schema.Resource {
	return &schema.Resource{
		Description: "Manages dashboard_panels.",

		CreateContext: resourceDashboardPanelCreate,
		ReadContext:   resourceDashboardPanelRead,
		UpdateContext: resourceDashboardPanelUpdate,
		DeleteContext: resourceDashboardPanelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"dashboard_id": {
				Description: "The id of the parent dashboard",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the dashboard_panel",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"y": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"w": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"h": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"params": {
				Description: "The params JSON of the dashboard_panel. See rootly API docs for schema.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"legend": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"groups": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"datasets": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collection": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"group_by": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"filter": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operation": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"rules": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"operation": &schema.Schema{
																Type:     schema.TypeString,
																Required: true,
															},
															"condition": &schema.Schema{
																Type:     schema.TypeString,
																Required: true,
															},
															"key": &schema.Schema{
																Type:     schema.TypeString,
																Required: true,
															},
															"value": &schema.Schema{
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
									"aggregate": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cumulative": &schema.Schema{
													Type:     schema.TypeBool,
													Required: true,
												},
												"key": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"operation": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func unflattenParams(params map[string]interface{}) []map[string]interface{} {
	unflattened_params := make([]map[string]interface{}, 1, 1)
	unflattened_params[0] = params
	datasets := params["datasets"].([]interface{})
	for _, dataset := range datasets {
		d := dataset.(map[string]interface{})
		aggregate := d["aggregate"].(map[string]interface{})
		if aggregate["key"] == nil {
			unflattened_aggregate := make([]interface{}, 0, 1)
			d["aggregate"] = unflattened_aggregate
		} else {
			unflattened_aggregate := make([]interface{}, 1, 1)
			unflattened_aggregate[0] = d["aggregate"]
			d["aggregate"] = unflattened_aggregate
		}
	}
	unflattened_params[0]["datasets"] = datasets
	unflattened_params[0]["legend"] = []interface{}{params["legend"]}
	return unflattened_params
}

func flattenParams(params []interface{}) map[string]interface{} {
	first_params := params[0].(map[string]interface{})
	datasets := first_params["datasets"].([]interface{})
	flattened_params := make(map[string]interface{})
	flattened_params["display"] = first_params["display"]
	flattened_params["legend"] = first_params["legend"].([]interface{})[0]
	flattened_datasets := make([]interface{}, len(datasets), len(datasets))
	for i, dataset := range datasets {
		d := dataset.(map[string]interface{})
		flattened_dataset := make(map[string]interface{})
		flattened_dataset["collection"] = d["collection"]
		flattened_dataset["name"] = d["name"]
		flattened_dataset["group_by"] = d["group_by"]
		flattened_dataset["filter"] = d["filter"]
		aggregate := d["aggregate"].([]interface{})
		if len(aggregate) == 1 {
			flattened_dataset["aggregate"] = aggregate[0]
		} else {
			flattened_dataset["aggregate"] = make(map[string]interface{})
		}
		flattened_datasets[i] = flattened_dataset
	}
	flattened_params["datasets"] = flattened_datasets
	return flattened_params
}

func resourceDashboardPanelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	dashboard_id := d.Get("dashboard_id").(string)
	name := d.Get("name").(string)

	tflog.Trace(ctx, fmt.Sprintf("Creating DashboardPanel: %s", name))

	s := &client.DashboardPanel{
		Name: name,
	}

	if value, ok := d.GetOk("position"); ok {
		s.Position = value.([]interface{})[0].(map[string]interface{})
	}

	if value, ok := d.GetOk("params"); ok {
		s.Params = flattenParams(value.([]interface{}))
	}

	res, err := c.CreateDashboardPanel(dashboard_id, s)
	if err != nil {
		return diag.Errorf("Error creating dashboard_panel: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a dashboard_panel resource: %v (%s)", name, d.Id()))

	return resourceDashboardPanelRead(ctx, d, meta)
}

func resourceDashboardPanelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading DashboardPanel: %s", d.Id()))

	dashboard_panel, err := c.GetDashboardPanel(d.Id(), new(rootlygo.GetDashboardPanelParams))
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("DashboardPanel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading dashboard_panel: %s", d.Id())
	}

	d.Set("name", dashboard_panel.Name)
	d.Set("params", unflattenParams(dashboard_panel.Params))
	d.Set("position", []map[string]interface{}{dashboard_panel.Position})

	return nil
}

func resourceDashboardPanelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating DashboardPanel: %s", d.Id()))

	name := d.Get("name").(string)

	s := &client.DashboardPanel{
		Name: name,
	}

	if value, ok := d.GetOk("position"); ok {
		s.Position = value.([]interface{})[0].(map[string]interface{})
	}

	if value, ok := d.GetOk("params"); ok {
		s.Params = flattenParams(value.([]interface{}))
	}

	_, err := c.UpdateDashboardPanel(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating dashboard_panel: %s", err.Error())
	}

	return resourceDashboardPanelRead(ctx, d, meta)
}

func resourceDashboardPanelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting DashboardPanel: %s", d.Id()))

	err := c.DeleteDashboardPanel(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("DashboardPanel (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting dashboard_panel: %s", err.Error())
	}

	d.SetId("")

	return nil
}
