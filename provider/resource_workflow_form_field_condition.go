package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	"github.com/rootlyhq/terraform-provider-rootly/tools"
)

func resourceWorkflowFormFieldCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowFormFieldConditionCreate,
		ReadContext:   resourceWorkflowFormFieldConditionRead,
		UpdateContext: resourceWorkflowFormFieldConditionUpdate,
		DeleteContext: resourceWorkflowFormFieldConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"workflow_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The workflow for this condition",
			},

			"form_field_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The custom field for this condition",
			},

			"incident_condition": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "ANY",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The trigger condition. Value must be one of `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.",
			},

			"values": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"selected_option_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},

			"selected_user_ids": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				DiffSuppressFunc: tools.EqualIgnoringOrder,
				Computed:         true,
				Required:         false,
				Optional:         true,
				Description:      "",
			},
		},
	}
}

func resourceWorkflowFormFieldConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating WorkflowFormFieldCondition"))

	s := &client.WorkflowFormFieldCondition{}

	if value, ok := d.GetOkExists("workflow_id"); ok {
		s.WorkflowId = value.(string)
	}
	if value, ok := d.GetOkExists("form_field_id"); ok {
		s.FormFieldId = value.(string)
	}
	if value, ok := d.GetOkExists("incident_condition"); ok {
		s.IncidentCondition = value.(string)
	}
	if value, ok := d.GetOkExists("values"); ok {
		s.Values = value.([]interface{})
	}
	if value, ok := d.GetOkExists("selected_option_ids"); ok {
		s.SelectedOptionIds = value.([]interface{})
	}
	if value, ok := d.GetOkExists("selected_user_ids"); ok {
		s.SelectedUserIds = value.([]interface{})
	}

	res, err := c.CreateWorkflowFormFieldCondition(s)
	if err != nil {
		return diag.Errorf("Error creating workflow_form_field_condition: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a workflow_form_field_condition resource: %s", d.Id()))

	return resourceWorkflowFormFieldConditionRead(ctx, d, meta)
}

func resourceWorkflowFormFieldConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading WorkflowFormFieldCondition: %s", d.Id()))

	item, err := c.GetWorkflowFormFieldCondition(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowFormFieldCondition (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading workflow_form_field_condition: %s", d.Id())
	}

	d.Set("workflow_id", item.WorkflowId)
	d.Set("form_field_id", item.FormFieldId)
	d.Set("incident_condition", item.IncidentCondition)
	d.Set("values", item.Values)
	d.Set("selected_option_ids", item.SelectedOptionIds)
	d.Set("selected_user_ids", item.SelectedUserIds)

	return nil
}

func resourceWorkflowFormFieldConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating WorkflowFormFieldCondition: %s", d.Id()))

	s := &client.WorkflowFormFieldCondition{}

	if d.HasChange("workflow_id") {
		s.WorkflowId = d.Get("workflow_id").(string)
	}
	if d.HasChange("form_field_id") {
		s.FormFieldId = d.Get("form_field_id").(string)
	}
	if d.HasChange("incident_condition") {
		s.IncidentCondition = d.Get("incident_condition").(string)
	}
	if d.HasChange("values") {
		s.Values = d.Get("values").([]interface{})
	}
	if d.HasChange("selected_option_ids") {
		s.SelectedOptionIds = d.Get("selected_option_ids").([]interface{})
	}
	if d.HasChange("selected_user_ids") {
		s.SelectedUserIds = d.Get("selected_user_ids").([]interface{})
	}

	_, err := c.UpdateWorkflowFormFieldCondition(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating workflow_form_field_condition: %s", err.Error())
	}

	return resourceWorkflowFormFieldConditionRead(ctx, d, meta)
}

func resourceWorkflowFormFieldConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting WorkflowFormFieldCondition: %s", d.Id()))

	err := c.DeleteWorkflowFormFieldCondition(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("WorkflowFormFieldCondition (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting workflow_form_field_condition: %s", err.Error())
	}

	d.SetId("")

	return nil
}
