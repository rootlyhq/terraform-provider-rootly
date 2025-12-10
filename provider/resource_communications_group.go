package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/diffsuppressfunc"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

// filterPropertiesName removes the deprecated 'name' field from communication_group_conditions properties
func filterPropertiesName(conditions []interface{}) []interface{} {
	if conditions == nil {
		return nil
	}

	filtered := make([]interface{}, len(conditions))
	for i, c := range conditions {
		if conditionMap, ok := c.(map[string]interface{}); ok {
			// Create a copy of the condition map
			filteredCondition := make(map[string]interface{})
			for k, v := range conditionMap {
				filteredCondition[k] = v
			}

			// Filter the properties array
			if props, ok := conditionMap["properties"].([]interface{}); ok {
				filteredProps := make([]interface{}, len(props))
				for j, p := range props {
					if propMap, ok := p.(map[string]interface{}); ok {
						// Create a copy without the 'name' field
						filteredProp := make(map[string]interface{})
						for k, v := range propMap {
							if k != "name" {
								filteredProp[k] = v
							}
						}
						filteredProps[j] = filteredProp
					} else {
						filteredProps[j] = p
					}
				}
				filteredCondition["properties"] = filteredProps
			}

			filtered[i] = filteredCondition
		} else {
			filtered[i] = c
		}
	}

	return filtered
}

func resourceCommunicationsGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommunicationsGroupCreate,
		ReadContext:   resourceCommunicationsGroupRead,
		UpdateContext: resourceCommunicationsGroupUpdate,
		DeleteContext: resourceCommunicationsGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			emailChannel, emailSet := d.GetOk("email_channel")
			smsChannel, smsSet := d.GetOk("sms_channel")

			// Check if at least one channel is explicitly set to true
			hasEmailChannel := emailSet && emailChannel.(bool)
			hasSmsChannel := smsSet && smsChannel.(bool)

			if !hasEmailChannel && !hasSmsChannel {
				return fmt.Errorf("at least one of 'email_channel' or 'sms_channel' must be set to true")
			}

			return nil
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The name of the communications group",
			},

			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The slug of the communications group",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The description of the communications group",
			},

			"communication_type_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "The communication type ID",
			},

			"is_private": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Whether the group is private. Value must be one of true or false",
			},

			"condition_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "any",
				Required:     false,
				Optional:     true,
				Sensitive:    false,
				ForceNew:     false,
				WriteOnly:    false,
				Description:  "Condition type. Value must be one of `any`, `all`.",
				ValidateFunc: validation.StringInSlice([]string{"any", "all"}, false),
			},

			"sms_channel": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "SMS channel enabled. Value must be one of true or false",
			},

			"email_channel": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    false,
				WriteOnly:   false,
				Description: "Email channel enabled. Value must be one of true or false",
			},

			"communication_group_conditions": &schema.Schema{
				Type:             schema.TypeList,
				Computed:         false,
				Required:         false,
				Optional:         true,
				Sensitive:        false,
				ForceNew:         false,
				WriteOnly:        false,
				Description:      "Group conditions",
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "ID of the condition",
						},

						"condition": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "Condition",
						},

						"property_type": &schema.Schema{
							Type:         schema.TypeString,
							Default:      "service",
							Required:     false,
							Optional:     true,
							Sensitive:    false,
							ForceNew:     false,
							WriteOnly:    false,
							Description:  "Property type. Value must be one of `service`, `severity`, `functionality`, `group`, `incident_type`.",
							ValidateFunc: validation.StringInSlice([]string{"service", "severity", "functionality", "group", "incident_type"}, false),
						},

						"properties": &schema.Schema{
							Type:             schema.TypeList,
							Computed:         false,
							Required:         false,
							Optional:         true,
							Sensitive:        false,
							ForceNew:         false,
							WriteOnly:        false,
							Description:      "Properties",
							DiffSuppressFunc: tools.EqualIgnoringOrder,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": &schema.Schema{
										Type:             schema.TypeString,
										Computed:         true,
										Required:         false,
										Optional:         true,
										Sensitive:        false,
										ForceNew:         false,
										WriteOnly:        false,
										Description:      "",
										Deprecated:       "This field is deprecated and will be removed in a future version",
										DiffSuppressFunc: diffsuppressfunc.Skip,
									},

									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    true,
										Sensitive:   false,
										ForceNew:    false,
										WriteOnly:   false,
										Description: "",
									},
								},
							},
						},
					},
				},
			},

			"communication_group_members": &schema.Schema{
				Type:             schema.TypeList,
				Computed:         false,
				Required:         false,
				Optional:         true,
				Sensitive:        false,
				ForceNew:         false,
				WriteOnly:        false,
				Description:      "Group members",
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "ID of the group member",
						},

						"user_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "User ID",
						},
					},
				},
			},

			"communication_external_group_members": &schema.Schema{
				Type:             schema.TypeList,
				Computed:         false,
				Required:         false,
				Optional:         true,
				Sensitive:        false,
				ForceNew:         false,
				WriteOnly:        false,
				Description:      "External group members",
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "ID of the external group member",
						},

						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "Name of the external member",
						},

						"email": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "Email of the external member",
						},

						"phone_number": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    true,
							Sensitive:   false,
							ForceNew:    false,
							WriteOnly:   false,
							Description: "Phone number of the external member",
						},
					},
				},
			},
		},
	}
}

func resourceCommunicationsGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating CommunicationsGroup"))

	s := &client.CommunicationsGroup{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("slug"); ok {
		s.Slug = value.(string)
	}
	if value, ok := d.GetOkExists("description"); ok {
		s.Description = value.(string)
	}
	if value, ok := d.GetOkExists("communication_type_id"); ok {
		s.CommunicationTypeId = value.(string)
	}
	if value, ok := d.GetOkExists("is_private"); ok {
		s.IsPrivate = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("condition_type"); ok {
		s.ConditionType = value.(string)
	}
	if value, ok := d.GetOkExists("sms_channel"); ok {
		s.SmsChannel = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("email_channel"); ok {
		s.EmailChannel = tools.Bool(value.(bool))
	}
	if value, ok := d.GetOkExists("communication_group_conditions"); ok {
		s.CommunicationGroupConditions = filterPropertiesName(value.([]interface{}))
	}
	if value, ok := d.GetOkExists("communication_group_members"); ok {
		s.CommunicationGroupMembers = value.([]interface{})
	}
	if value, ok := d.GetOkExists("communication_external_group_members"); ok {
		s.CommunicationExternalGroupMembers = value.([]interface{})
	}

	res, err := c.CreateCommunicationsGroup(s)
	if err != nil {
		return diag.Errorf("Error creating communications_group: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a communications_group resource: %s", d.Id()))

	return resourceCommunicationsGroupRead(ctx, d, meta)
}

func resourceCommunicationsGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading CommunicationsGroup: %s", d.Id()))

	item, err := c.GetCommunicationsGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CommunicationsGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading communications_group: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("slug", item.Slug)
	d.Set("description", item.Description)
	d.Set("communication_type_id", item.CommunicationTypeId)
	d.Set("is_private", item.IsPrivate)
	d.Set("condition_type", item.ConditionType)
	d.Set("sms_channel", item.SmsChannel)
	d.Set("email_channel", item.EmailChannel)

	if item.CommunicationGroupConditions != nil {
		processed_items_communication_group_conditions := make([]map[string]interface{}, 0)

		for _, c := range item.CommunicationGroupConditions {
			if rawItem, ok := c.(map[string]interface{}); ok {
				// Create a new map with only the fields defined in the schema
				processed_item_communication_group_conditions := map[string]interface{}{
					"id":            rawItem["id"],
					"condition":     rawItem["condition"],
					"property_type": rawItem["property_type"],
					"properties":    rawItem["properties"],
				}
				processed_items_communication_group_conditions = append(processed_items_communication_group_conditions, processed_item_communication_group_conditions)
			}
		}

		d.Set("communication_group_conditions", processed_items_communication_group_conditions)
	} else {
		d.Set("communication_group_conditions", nil)
	}

	if item.CommunicationGroupMembers != nil {
		processed_items_communication_group_members := make([]map[string]interface{}, 0)

		for _, c := range item.CommunicationGroupMembers {
			if rawItem, ok := c.(map[string]interface{}); ok {
				// Create a new map with only the fields defined in the schema
				processed_item_communication_group_members := map[string]interface{}{
					"id":      rawItem["id"],
					"user_id": rawItem["user_id"],
				}
				processed_items_communication_group_members = append(processed_items_communication_group_members, processed_item_communication_group_members)
			}
		}

		d.Set("communication_group_members", processed_items_communication_group_members)
	} else {
		d.Set("communication_group_members", nil)
	}

	if item.CommunicationExternalGroupMembers != nil {
		processed_items_communication_external_group_members := make([]map[string]interface{}, 0)

		for _, c := range item.CommunicationExternalGroupMembers {
			if rawItem, ok := c.(map[string]interface{}); ok {
				// Create a new map with only the fields defined in the schema
				processed_item_communication_external_group_members := map[string]interface{}{
					"id":           rawItem["id"],
					"name":         rawItem["name"],
					"email":        rawItem["email"],
					"phone_number": rawItem["phone_number"],
				}
				processed_items_communication_external_group_members = append(processed_items_communication_external_group_members, processed_item_communication_external_group_members)
			}
		}

		d.Set("communication_external_group_members", processed_items_communication_external_group_members)
	} else {
		d.Set("communication_external_group_members", nil)
	}

	return nil
}

func resourceCommunicationsGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating CommunicationsGroup: %s", d.Id()))

	s := &client.CommunicationsGroup{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("slug") {
		s.Slug = d.Get("slug").(string)
	}
	if d.HasChange("description") {
		s.Description = d.Get("description").(string)
	}
	if d.HasChange("communication_type_id") {
		s.CommunicationTypeId = d.Get("communication_type_id").(string)
	}
	if d.HasChange("is_private") {
		s.IsPrivate = tools.Bool(d.Get("is_private").(bool))
	}
	if d.HasChange("condition_type") {
		s.ConditionType = d.Get("condition_type").(string)
	}
	if d.HasChange("sms_channel") {
		s.SmsChannel = tools.Bool(d.Get("sms_channel").(bool))
	}
	if d.HasChange("email_channel") {
		s.EmailChannel = tools.Bool(d.Get("email_channel").(bool))
	}

	if d.HasChange("communication_group_conditions") {
		if value, ok := d.GetOk("communication_group_conditions"); value != nil && ok {
			s.CommunicationGroupConditions = filterPropertiesName(value.([]interface{}))
		} else {
			s.CommunicationGroupConditions = []interface{}{}
		}
	}

	if d.HasChange("communication_group_members") {
		if value, ok := d.GetOk("communication_group_members"); value != nil && ok {
			s.CommunicationGroupMembers = value.([]interface{})
		} else {
			s.CommunicationGroupMembers = []interface{}{}
		}
	}

	if d.HasChange("communication_external_group_members") {
		if value, ok := d.GetOk("communication_external_group_members"); value != nil && ok {
			s.CommunicationExternalGroupMembers = value.([]interface{})
		} else {
			s.CommunicationExternalGroupMembers = []interface{}{}
		}
	}

	_, err := c.UpdateCommunicationsGroup(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating communications_group: %s", err.Error())
	}

	return resourceCommunicationsGroupRead(ctx, d, meta)
}

func resourceCommunicationsGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting CommunicationsGroup: %s", d.Id()))

	err := c.DeleteCommunicationsGroup(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("CommunicationsGroup (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting communications_group: %s", err.Error())
	}

	d.SetId("")

	return nil
}
