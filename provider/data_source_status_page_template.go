package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSourceStatusPageTemplate() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceStatusPageTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			
			"status_page_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"title": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"body": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"update_status": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"should_notify_subscribers": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

				"enabled": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				

			"position": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			

			"created_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"updated_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			
		},
	}
}

func dataSourceStatusPageTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListStatusPageTemplatesParams)
	page_size := 1
	params.PageSize = &page_size

	

	status_page_id := d.Get("status_page_id").(string)
			items, err := c.ListStatusPageTemplates(status_page_id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("status_page_template not found")
	}
	item, _ := items[0].(*client.StatusPageTemplate)

	d.SetId(item.ID)

	return nil
}
