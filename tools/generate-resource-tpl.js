const inflect = require('inflect');

module.exports = (name, resourceSchema, collectionSchema, pathIdField, createResourceSchema) => {
	const namePlural = inflect.pluralize(name)
	const nameCamel = inflect.camelize(name)
	const nameCamelPlural = inflect.camelize(namePlural)

return `package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resource${nameCamel}() *schema.Resource{
	return &schema.Resource{
		CreateContext: resource${nameCamel}Create,
		ReadContext: resource${nameCamel}Read,
		UpdateContext: resource${nameCamel}Update,
		DeleteContext: resource${nameCamel}Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			${schemaFields(resourceSchema, createResourceSchema)}
		},
	}
}

func resource${nameCamel}Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating ${nameCamel}"))

	s := &client.${nameCamel}{}

	${createResourceFields(name, resourceSchema)}

	res, err := c.Create${nameCamel}(s)
	if err != nil {
		return diag.Errorf("Error creating ${name}: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a ${name} resource: %s", d.Id()))

	return resource${nameCamel}Read(ctx, d, meta)
}

func resource${nameCamel}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading ${nameCamel}: %s", d.Id()))

	item, err := c.Get${nameCamel}(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("${nameCamel} (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading ${name}: %s", d.Id())
	}

	${setResourceFields(name, resourceSchema)}

	return nil
}

func resource${nameCamel}Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating ${nameCamel}: %s", d.Id()))

	s := &client.${nameCamel}{}

	${updateResourceFields(name, resourceSchema)}

	_, err := c.Update${nameCamel}(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating ${name}: %s", err.Error())
	}

	return resource${nameCamel}Read(ctx, d, meta)
}

func resource${nameCamel}Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting ${nameCamel}: %s", d.Id()))

	err := c.Delete${nameCamel}(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("${nameCamel} (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting ${name}: %s", err.Error())
	}

	d.SetId("")

	return nil
}
`}

function excludeDateFields(field) {
	return field !== 'created_at' && field !== 'updated_at'
}

function setResourceFields(name, resourceSchema) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		return `d.Set("${field}", item.${inflect.camelize(field)})`
	}).join('\n  ')
}

function createResourceFields(name, resourceSchema) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		const schema = resourceSchema.properties[field]
		return`  if value, ok := d.GetOkExists("${field}"); ok {
		s.${inflect.camelize(field)} = value.(${jsonapiToGoType(schema.type)})
	}`
	}).join('\n  ')
}

function updateResourceFields(name, resourceSchema) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		const schema = resourceSchema.properties[field]
		return`  if d.HasChange("${field}") {
		s.${inflect.camelize(field)} = d.Get("${field}").(${jsonapiToGoType(schema.type)})
	}`
	}).join('\n  ')
}

function jsonapiToGoType(type) {
	switch (type) {
		case 'string':
			return 'string'
		case 'number':
			return 'int'
		case 'boolean':
			return 'bool'
		case 'array':
			return '[]interface{}'
		case 'object':
			return 'interface{}'
		default:
			return 'string'
	}
}

function schemaFields(resourceSchema, createResourceSchema) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		return schemaField(field, resourceSchema, createResourceSchema)
	}).join('\n')
}

function schemaField(name, resourceSchema, createResourceSchema) {
	const schema = resourceSchema.properties[name]
	const optional = (createResourceSchema.required || []).indexOf(name) === -1 ? "true" : "false"
	const required = (createResourceSchema.required || []).indexOf(name) === -1 ? "false" : "true"
	const description = (schema.description || '').replace(/"/g, '\\"')
	switch (schema.type) {
		case 'string':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeString,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`
		case 'number':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeInt,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`
		case 'boolean':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeBool,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`
		case 'array':
			if (schema.items && schema.items.type === "object") {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Computed: ${optional},
					Required: ${required},
					Optional: ${optional},
					Description: "${description}",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
							"name": &schema.Schema{
								Type: schema.TypeString,
								Required: true,
							},
						},
					},
				},
				`
			} else {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: ${optional},
					Required: ${required},
					Optional: ${optional},
					Description: "${description}",
				},
				`
			}
		case 'object':
		default:
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				Description: "${description}",
			},
			`
	}
}
