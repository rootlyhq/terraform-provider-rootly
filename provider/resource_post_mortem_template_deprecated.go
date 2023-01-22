package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePostmortemTemplateDeprecated() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostmortemTemplateDeprecatedCreate,
		ReadContext:   resourcePostmortemTemplateDeprecatedRead,
		UpdateContext: resourcePostmortemTemplateDeprecatedUpdate,
		DeleteContext: resourcePostmortemTemplateDeprecatedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the postmortem template",
			},

			"default": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: "Default selected template when editing a postmortem",
			},

			"content": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "The postmortem template. Liquid syntax and markdown are supported.",

				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return len(old) != 0
				},
			},
		},
	}
}

func resourcePostmortemTemplateDeprecatedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Rename to rootly_post_mortem_template (note underscores). Name changed to match Rootly API.")
}

func resourcePostmortemTemplateDeprecatedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Rename to rootly_post_mortem_template (note underscores). Name changed to match Rootly API.")
}

func resourcePostmortemTemplateDeprecatedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Rename to rootly_post_mortem_template (note underscores). Name changed to match Rootly API.")
}

func resourcePostmortemTemplateDeprecatedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Rename to rootly_post_mortem_template (note underscores). Name changed to match Rootly API.")
}
