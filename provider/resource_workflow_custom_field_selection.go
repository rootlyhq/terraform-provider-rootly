package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceWorkflowCustomFieldSelection() *schema.Resource {
	return &schema.Resource{
		Description: "Manages workflow custom field selections (triggers).",

		CreateContext: resourceWorkflowCustomFieldSelectionCreate,
		ReadContext:   resourceWorkflowCustomFieldSelectionRead,
		UpdateContext: resourceWorkflowCustomFieldSelectionUpdate,
		DeleteContext: resourceWorkflowCustomFieldSelectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"workflow_id": {
				Description:  "The ID of the workflow",
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"custom_field_id": {
				Description:  "The ID of the custom field",
				Type: schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"incident_condition": {
				Description: "The trigger condition",
				Type: schema.TypeString,
				Optional: true,
				Default: "ANY",
				ValidateFunc: validation.StringInSlice([]string{
					"IS",
					"ANY",
					"CONTAINS",
					"CONTAINS_ALL",
					"NONE",
					"SET",
					"UNSET",
				}, false),
			},
			"values": {
				Description:  "Custom field values to associate with this custom field trigger",
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"selected_option_ids": {
				Description:  "Custom field options to associate with this custom field trigger",
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceWorkflowCustomFieldSelectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	workflowId := d.Get("workflow_id").(string)
	customFieldId := d.Get("custom_field_id").(int)

	tflog.Trace(ctx, fmt.Sprintf("Creating workflow custom field selection"))

	s := &client.WorkflowCustomFieldSelection{
		CustomFieldId: customFieldId,
	}

	if value, ok := d.GetOk("incident_condition"); ok {
		s.IncidentCondition = value.(string)
	}

	if value, ok := d.GetOk("values"); ok {
		s.Values = value.([]interface{})
	}

	if value, ok := d.GetOk("selected_option_ids"); ok {
		s.SelectedOptionIds = value.([]interface{})
	}

	res, err := c.CreateWorkflowCustomFieldSelection(workflowId, s)
	if err != nil {
		return diag.Errorf("Error creating workflow custom field selection: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created an workflow custom field selection resource: (%s)", d.Id()))

	return resourceWorkflowCustomFieldSelectionRead(ctx, d, meta)
}

func resourceWorkflowCustomFieldSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading workflow custom field selection: %s", d.Id()))

	res, err := c.GetWorkflowCustomFieldSelection(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowCustomFieldSelection (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow custom field selection: %s", d.Id())
	}

	d.Set("custom_field_id", res.CustomFieldId)
	d.Set("workflow_id", res.WorkflowId)
	d.Set("incident_condition", res.IncidentCondition)
	d.Set("values", res.Values)
	d.Set("selected_option_ids", res.SelectedOptionIds)

	return nil
}

func resourceWorkflowCustomFieldSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating workflow custom field selection: %s", d.Id()))

	incident_condition := d.Get("incident_condition").(string)

	s := &client.WorkflowCustomFieldSelection{
		IncidentCondition: incident_condition,
	}

	if d.HasChange("values") {
		s.Values = d.Get("values").([]interface{})
	}

	if d.HasChange("selected_option_ids") {
		s.SelectedOptionIds = d.Get("selected_option_ids").([]interface{})
	}

	_, err := c.UpdateWorkflowCustomFieldSelection(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow custom field selection: %s", err.Error())
	}

	return resourceWorkflowCustomFieldSelectionRead(ctx, d, meta)
}

func resourceWorkflowCustomFieldSelectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting workflow custom field selection: %s", d.Id()))

	err := c.DeleteWorkflowCustomFieldSelection(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowCustomFieldSelection (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow custom field selection: %s", err.Error())
	}

	d.SetId("")

	return nil
}
