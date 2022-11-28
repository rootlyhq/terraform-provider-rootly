package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceWebhooksDelivery() *schema.Resource{
	return &schema.Resource{
		
		ReadContext: resourceWebhooksDeliveryRead,
		
		
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"endpoint_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			

			"delivered": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "",
			},
			

			"attempts": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			

			"payload": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			

			"last_attempt_at": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "",
			},
			
		},
	}
}


func resourceWebhooksDeliveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WebhooksDelivery: %s", d.Id()))

	item, err := c.GetWebhooksDelivery(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WebhooksDelivery (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading webhooks_delivery: %s", d.Id())
	}

	d.Set("endpoint_id", item.EndpointId)
  d.Set("delivered", item.Delivered)
  d.Set("attempts", item.Attempts)
  d.Set("payload", item.Payload)
  d.Set("last_attempt_at", item.LastAttemptAt)

	return nil
}

