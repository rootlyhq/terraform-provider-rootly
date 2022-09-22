package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceIncidentRole() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceIncidentRoleRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"slug": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"created_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceIncidentRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentRolesParams)
	page_size := 1
	params.PageSize = &page_size

	
				slug := d.Get("slug").(string)
				params.FilterSlug = &slug
			

				name := d.Get("name").(string)
				params.FilterName = &name
			

	items, err := c.ListIncidentRoles(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("incident_role not found")
	}
	item, _ := items[0].(*client.IncidentRole)

	d.SetId(item.ID)

	return nil
}
