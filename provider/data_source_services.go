package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/rootly-go"
)

func dataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServicesRead,
		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"backstage_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"pagerduty_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"opsgenie_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":                         {Type: schema.TypeString, Computed: true},
						"name":                       {Type: schema.TypeString, Computed: true},
						"slug":                       {Type: schema.TypeString, Computed: true},
						"description":                {Type: schema.TypeString, Computed: true},
						"public_description":         {Type: schema.TypeString, Computed: true},
						"color":                      {Type: schema.TypeString, Computed: true},
						"position":                   {Type: schema.TypeInt, Computed: true},
						"backstage_id":               {Type: schema.TypeString, Computed: true},
						"external_id":                {Type: schema.TypeString, Computed: true},
						"pagerduty_id":               {Type: schema.TypeString, Computed: true},
						"opsgenie_id":                {Type: schema.TypeString, Computed: true},
						"cortex_id":                  {Type: schema.TypeString, Computed: true},
						"service_now_ci_sys_id":      {Type: schema.TypeString, Computed: true},
						"github_repository_name":     {Type: schema.TypeString, Computed: true},
						"github_repository_branch":   {Type: schema.TypeString, Computed: true},
						"gitlab_repository_name":     {Type: schema.TypeString, Computed: true},
						"gitlab_repository_branch":   {Type: schema.TypeString, Computed: true},
						"kubernetes_deployment_name": {Type: schema.TypeString, Computed: true},
						"notify_emails":              {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"environment_ids":            {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"service_ids":                {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"owner_group_ids":            {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"owner_user_ids":             {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_urgency_id":           {Type: schema.TypeString, Computed: true},
						"escalation_policy_id":       {Type: schema.TypeString, Computed: true},
						"alerts_email_enabled":       {Type: schema.TypeBool, Computed: true},
						"alerts_email_address":       {Type: schema.TypeString, Computed: true},
						"slack_channels":             {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"slack_aliases":              {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_broadcast_enabled":    {Type: schema.TypeBool, Computed: true},
						"alert_broadcast_channel":    {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"incident_broadcast_enabled": {Type: schema.TypeBool, Computed: true},
						"incident_broadcast_channel": {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"properties":                 {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					},
				},
			},
		},
	}
}

func dataSourceServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListServicesParams)

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	params.FilterSlug = &slug
	params.FilterName = &name

	services, err := c.ListServices(params)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceSchema := dataSourceServices().Schema["services"].Elem.(*schema.Resource).Schema
	tf_services := make([]interface{}, len(services), len(services))
	for i, service := range services {
		c, _ := service.(*client.Service)
		tf_services[i] = filterMapKeys(structToLowerFirstMap(*c), serviceSchema)
	}

	if err := d.Set("services", tf_services); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
