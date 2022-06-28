package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_host": {
					Description: "The Rootly API host. Defaults to https://api.rootly.com",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_URL", "https://api.rootly.com"),
				},
				"api_token": {
					Description: "The Rootly API Token. Generate it from your account at https://rootly.com/account",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_TOKEN", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"rootly_cause":         resourceCause(),
				"rootly_environment":   resourceEnvironment(),
				"rootly_functionality": resourceFunctionality(),
				"rootly_incident_role": resourceIncidentRole(),
				"rootly_incident_type": resourceIncidentType(),
				"rootly_service":       resourceService(),
				"rootly_severity":      resourceSeverity(),
				"rootly_team":          resourceTeam(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		host := d.Get("api_host").(string)
		token := d.Get("api_token").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		cli, err := client.NewClient(host, token, RootlyUserAgent(version))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Rootly client",
				Detail:   "Unable to authenticate user for authenticated Rootly client",
			})

			return nil, diags
		}

		return cli, diags
	}
}

func RootlyUserAgent(version string) string {
	return fmt.Sprintf("Rootly Terraform Provider/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", version, meta.SDKVersionString())
}
