package provider

import (
	"context"
	"strconv"
	"time"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceIncidentRoles() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceIncidentRolesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"incident_roles": &schema.Schema{
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type: schema.TypeString,
							Computed: true,
						},
						"slug": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"summary": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIncidentRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentRolesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)

	params.FilterSlug = &slug
	params.FilterName = &name
	params.FilterEnabled = &enabled

	incident_roles, err := c.ListIncidentRoles(params)
	if err != nil {
		return diag.FromErr(err)
	}

	tf_incident_roles := make([]interface{}, len(incident_roles), len(incident_roles))
	for i, incident_role := range incident_roles {
		c, _ := incident_role.(*client.IncidentRole)
		tf_incident_roles[i] = structToLowerFirstMap(*c)
	}

	if err := d.Set("incident_roles", tf_incident_roles); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
