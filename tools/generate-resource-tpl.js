const inflect = require('./inflect');

function includeTools(resourceSchema) {
	for (var key in resourceSchema.properties) {
		if (resourceSchema.properties[key].type === "boolean") {
			return true;
		}
	}
	return false;
}

module.exports = (name, resourceSchema, requiredFields, pathIdField, exclude = []) => {
	const nameCamel = inflect.camelize(name)
	const tools = includeTools(resourceSchema) && exclude.length === 0 ? `"github.com/rootlyhq/terraform-provider-rootly/tools"` : ""

definition = `package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	${tools}
)

func resource${nameCamel}() *schema.Resource{
	return &schema.Resource{
		${(!exclude.includes("create")) ? `CreateContext: resource${nameCamel}Create,` : ""}
		${(!exclude.includes("read")) ? `ReadContext: resource${nameCamel}Read,` : ""}
		${(!exclude.includes("update")) ? `UpdateContext: resource${nameCamel}Update,` : ""}
		${(!exclude.includes("delete")) ? `DeleteContext: resource${nameCamel}Delete,` : ""}
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			${schemaFields(resourceSchema, requiredFields, pathIdField)}
		},
	}
}

`

	if (!exclude.includes("create")) {

		definition = definition + `
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

`
	}

	if (!exclude.includes("read")) {

		definition = definition + `
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

`
	}

	if (!exclude.includes("update")) {

		definition = definition + `
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

`
	}

	if (!exclude.includes("delete")) {

		definition = definition + `
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

`
	}
	return definition;
}

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
		if (schema.type === "boolean") {
			return`  if value, ok := d.GetOkExists("${field}"); ok {
				s.${inflect.camelize(field)} = tools.Bool(value.(${jsonapiToGoType(schema.type)}))
			}`
		} else {
			return`  if value, ok := d.GetOkExists("${field}"); ok {
				s.${inflect.camelize(field)} = value.(${jsonapiToGoType(schema.type)})
			}`
		}
	}).join('\n  ')
}

function updateResourceFields(name, resourceSchema) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		const schema = resourceSchema.properties[field]
		if (schema.type === "boolean") {
			return`  if d.HasChange("${field}") {
				s.${inflect.camelize(field)} = tools.Bool(d.Get("${field}").(${jsonapiToGoType(schema.type)}))
			}`
		} else {
			return`  if d.HasChange("${field}") {
				s.${inflect.camelize(field)} = d.Get("${field}").(${jsonapiToGoType(schema.type)})
			}`
		}
	}).join('\n  ')
}

function jsonapiToGoType(type) {
	switch (type) {
		case 'string':
			return 'string'
		case 'integer':
		case 'number':
			return 'int'
		case 'boolean':
			return 'bool'
		case 'array':
			return '[]interface{}'
		case 'object':
			return 'map[string]interface{}'
		default:
			return 'interface{}'
	}
}

function schemaFields(resourceSchema, requiredFields, pathIdField) {
	return Object.keys(resourceSchema.properties).filter(excludeDateFields).map((field) => {
		return schemaField(field, resourceSchema, requiredFields, pathIdField)
	}).join('\n')
}

function annotatedDescription(schema) {
	const description = (schema.description || "").replace(/"/g, '\\"')
	if (schema.enum) {
		return `${description}. Value must be one of ${schema.enum.map((val) => `\`${val}\``).join(", ")}.`
	}
	if (schema.type === "object" && schema.properties.id && schema.properties.name) {
		return `Map must contain two fields, \`id\` and \`name\`. ${description}`
	}
	if (schema.type === "array" && schema.items && schema.items.enum) {
		return `${description}. Value must be one of ${schema.items.enum.map((val) => `\`${val}\``).join(", ")}.`
	}
	return description
}

function schemaField(name, resourceSchema, requiredFields, pathIdField) {
	const schema = resourceSchema.properties[name]
	const optional = (requiredFields || []).indexOf(name) === -1 || schema.enum ? "true" : "false"
	const required = (requiredFields || []).indexOf(name) === -1 || schema.enum ? "false" : "true"
	const description = annotatedDescription(schema)
	const forceNew = name === pathIdField ? "true" : "false"
	switch (schema.type) {
		case 'string':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeString,
				${schema.enum ? `Default: "${schema.enum[0]}",` : `Computed: ${optional},`}
				Required: ${required},
				Optional: ${optional},
				ForceNew: ${forceNew},
				Description: "${description}",
			},
			`
		case 'integer':
		case 'number':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeInt,
				Computed: ${optional},
				Required: ${required},
				Optional: ${optional},
				ForceNew: ${forceNew},
				Description: "${description}",
			},
			`
		case 'boolean':
			if (name === "enabled") {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				`
			}
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
			} else if (schema.items && schema.items.type === "string") {
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
			} else if (schema.items && schema.items.type === "number") {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
					Computed: ${optional},
					Required: ${required},
					Optional: ${optional},
					Description: "${description}",
				},
				`
			} else {
				console.log(`unsupported array field schema:`, name, schema)
				return ''
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
