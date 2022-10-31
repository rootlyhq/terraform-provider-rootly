package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	
)

func resourceIncidentFormFieldSelection() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceIncidentFormFieldSelectionCreate,
		ReadContext: resourceIncidentFormFieldSelectionRead,
		UpdateContext: resourceIncidentFormFieldSelectionUpdate,
		DeleteContext: resourceIncidentFormFieldSelectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			
			"incident_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: true,
				Description: "",
			},
			

			"form_field_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The custom field for this selection",
			},
			

			"value": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The selected value for text kind custom fields",
			},
			

				"selected_option_ids": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "",
				},
				

				"selected_user_ids": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
					Computed: true,
					Required: false,
					Optional: true,
					Description: "",
				},
				
		},
	}
}

func resourceIncidentFormFieldSelectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating IncidentFormFieldSelection"))

	s := &client.IncidentFormFieldSelection{}

	  if value, ok := d.GetOkExists("incident_id"); ok {
				s.IncidentId = value.(string)
			}
    if value, ok := d.GetOkExists("form_field_id"); ok {
				s.FormFieldId = value.(string)
			}
    if value, ok := d.GetOkExists("value"); ok {
				s.Value = value.(string)
			}
    if value, ok := d.GetOkExists("selected_option_ids"); ok {
				s.SelectedOptionIds = value.([]interface{})
			}
    if value, ok := d.GetOkExists("selected_user_ids"); ok {
				s.SelectedUserIds = value.([]interface{})
			}

	res, err := c.CreateIncidentFormFieldSelection(s)
	if err != nil {
		return diag.Errorf("Error creating incident_form_field_selection: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a incident_form_field_selection resource: %s", d.Id()))

	return resourceIncidentFormFieldSelectionRead(ctx, d, meta)
}

func resourceIncidentFormFieldSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading IncidentFormFieldSelection: %s", d.Id()))

	item, err := c.GetIncidentFormFieldSelection(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentFormFieldSelection (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading incident_form_field_selection: %s", d.Id())
	}

	d.Set("incident_id", item.IncidentId)
  d.Set("form_field_id", item.FormFieldId)
  d.Set("value", item.Value)
  d.Set("selected_option_ids", item.SelectedOptionIds)
  d.Set("selected_user_ids", item.SelectedUserIds)

	return nil
}

func resourceIncidentFormFieldSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating IncidentFormFieldSelection: %s", d.Id()))

	s := &client.IncidentFormFieldSelection{}

	  if d.HasChange("incident_id") {
				s.IncidentId = d.Get("incident_id").(string)
			}
    if d.HasChange("form_field_id") {
				s.FormFieldId = d.Get("form_field_id").(string)
			}
    if d.HasChange("value") {
				s.Value = d.Get("value").(string)
			}
    if d.HasChange("selected_option_ids") {
				s.SelectedOptionIds = d.Get("selected_option_ids").([]interface{})
			}
    if d.HasChange("selected_user_ids") {
				s.SelectedUserIds = d.Get("selected_user_ids").([]interface{})
			}

	_, err := c.UpdateIncidentFormFieldSelection(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating incident_form_field_selection: %s", err.Error())
	}

	return resourceIncidentFormFieldSelectionRead(ctx, d, meta)
}

func resourceIncidentFormFieldSelectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting IncidentFormFieldSelection: %s", d.Id()))

	err := c.DeleteIncidentFormFieldSelection(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("IncidentFormFieldSelection (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident_form_field_selection: %s", err.Error())
	}

	d.SetId("")

	return nil
}
