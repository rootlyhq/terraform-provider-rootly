// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
)

func dataSourceIncident() *schema.Resource {
	return &schema.Resource {
		ReadContext: dataSourceIncidentRead,
		Schema: map[string]*schema.Schema {
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
			},
			
			"kind": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"private": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			

			"status": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"user": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"severity": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"environments": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"services": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"functionalities": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

			"labels": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			

				"in_triage_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"started_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"detected_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"acknowledged_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"mitigated_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"resolved_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"closed_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"created_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				

				"updated_at": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				
		},
	}
}

func dataSourceIncidentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.ListIncidentsParams)
	page_size := 1
	params.PageSize = &page_size

	
				if value, ok := d.GetOkExists("kind"); ok {
					kind := value.(string)
					params.FilterKind = &kind
				}
			

				if value, ok := d.GetOkExists("status"); ok {
					status := value.(string)
					params.FilterStatus = &status
				}
			

				if value, ok := d.GetOkExists("private"); ok {
					private := value.(bool)
					params.FilterPrivate = &private
				}
			

				if value, ok := d.GetOkExists("severity"); ok {
					severity := value.(string)
					params.FilterSeverity = &severity
				}
			

				if value, ok := d.GetOkExists("labels"); ok {
					labels := value.(string)
					params.FilterLabels = &labels
				}
			

				if value, ok := d.GetOkExists("environments"); ok {
					environments := value.(string)
					params.FilterEnvironments = &environments
				}
			

				if value, ok := d.GetOkExists("functionalities"); ok {
					functionalities := value.(string)
					params.FilterFunctionalities = &functionalities
				}
			

				if value, ok := d.GetOkExists("services"); ok {
					services := value.(string)
					params.FilterServices = &services
				}
			

				created_at_gt := d.Get("created_at").(map[string]interface{})
				if value, exists := created_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterCreatedAtGt = &v
				}
			

				created_at_lt := d.Get("created_at").(map[string]interface{})
				if value, exists := created_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterCreatedAtLt = &v
				}
			

				updated_at_gt := d.Get("updated_at").(map[string]interface{})
				if value, exists := updated_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterUpdatedAtGt = &v
				}
			

				updated_at_lt := d.Get("updated_at").(map[string]interface{})
				if value, exists := updated_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterUpdatedAtLt = &v
				}
			

				started_at_gt := d.Get("started_at").(map[string]interface{})
				if value, exists := started_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterStartedAtGt = &v
				}
			

				started_at_lt := d.Get("started_at").(map[string]interface{})
				if value, exists := started_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterStartedAtLt = &v
				}
			

				detected_at_gt := d.Get("detected_at").(map[string]interface{})
				if value, exists := detected_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterDetectedAtGt = &v
				}
			

				detected_at_lt := d.Get("detected_at").(map[string]interface{})
				if value, exists := detected_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterDetectedAtLt = &v
				}
			

				acknowledged_at_gt := d.Get("acknowledged_at").(map[string]interface{})
				if value, exists := acknowledged_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterAcknowledgedAtGt = &v
				}
			

				acknowledged_at_lt := d.Get("acknowledged_at").(map[string]interface{})
				if value, exists := acknowledged_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterAcknowledgedAtLt = &v
				}
			

				mitigated_at_gt := d.Get("mitigated_at").(map[string]interface{})
				if value, exists := mitigated_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterMitigatedAtGt = &v
				}
			

				mitigated_at_lt := d.Get("mitigated_at").(map[string]interface{})
				if value, exists := mitigated_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterMitigatedAtLt = &v
				}
			

				resolved_at_gt := d.Get("resolved_at").(map[string]interface{})
				if value, exists := resolved_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterResolvedAtGt = &v
				}
			

				resolved_at_lt := d.Get("resolved_at").(map[string]interface{})
				if value, exists := resolved_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterResolvedAtLt = &v
				}
			

				closed_at_gt := d.Get("closed_at").(map[string]interface{})
				if value, exists := closed_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterClosedAtGt = &v
				}
			

				closed_at_lt := d.Get("closed_at").(map[string]interface{})
				if value, exists := closed_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterClosedAtLt = &v
				}
			

				in_triage_at_gt := d.Get("in_triage_at").(map[string]interface{})
				if value, exists := in_triage_at_gt["gt"]; exists {
					v := value.(string)
					params.FilterInTriageAtGt = &v
				}
			

				in_triage_at_lt := d.Get("in_triage_at").(map[string]interface{})
				if value, exists := in_triage_at_lt["lt"]; exists {
					v := value.(string)
					params.FilterInTriageAtLt = &v
				}
			

	items, err := c.ListIncidents(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("incident not found")
	}
	item, _ := items[0].(*client.Incident)

	d.SetId(item.ID)

	return nil
}
