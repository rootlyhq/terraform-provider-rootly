const inflect = require('inflect');

module.exports = (name, resourceSchema, filterParameters, pathIdField) => {
	const namePlural = inflect.pluralize(name)
	const nameCamel = inflect.camelize(name)
	const nameCamelPlural = inflect.camelize(namePlural)
	const strconvImport = pathIdField && resourceSchema.properties[pathIdField].type === 'number' ? `"strconv"` : ''

return `package provider

import (
	"context"
	${strconvImport}
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/schema"
)

func dataSource${nameCamel}() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSource${nameCamel}Read,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			${schemaFields(resourceSchema)}
		},
	}
}

func dataSource${nameCamel}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	params := new(rootlygo.List${nameCamelPlural}Params)
	page_size := 1
	params.PageSize = &page_size

	${setFilterFields(name, resourceSchema, filterParameters)}

	${listFn(nameCamelPlural, resourceSchema, pathIdField)}
	if err != nil {
		return diag.FromErr(err)
	}

	if len(items) == 0 {
		return diag.Errorf("${name} not found")
	}
	item, _ := items[0].(*client.${nameCamel})

	d.SetId(item.ID)

	return nil
}
`}

function listFn(nameCamelPlural, resourceSchema, pathIdField) {
	if (pathIdField) {
		if (resourceSchema.properties[pathIdField].type === 'number') {
			return `${pathIdField} := strconv.Itoa(d.Get("${pathIdField}").(int))
			items, err := c.List${nameCamelPlural}(${pathIdField}, params)`
		} else {
			return `${pathIdField} := d.Get("${pathIdField}").(string)
			items, err := c.List${nameCamelPlural}(${pathIdField}, params)`
		}
	} else {
		return `items, err := c.List${nameCamelPlural}(params)`
	}
}

function setFilterFields(name, resourceSchema, filterParameters) {
	return (filterParameters || []).filter((paramSchema) => {
		return paramSchema.name.match(/^filter/)
	}).map((paramSchema) => {
		const filterField = inflect.underscore(paramSchema.name.replace("filter[", "").replace("]", ""))
		const fieldSchema = resourceSchema.properties[filterField]
		if (fieldSchema) {
			return `
				${filterField} := d.Get("${filterField}").(${jsonapiToGoType(fieldSchema.type)})
				params.Filter${inflect.camelize(filterField)} = &${filterField}
			`
		}
	}).filter((x) => x).join('\n')
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
		default:
			return 'string'
	}
}

function schemaFields(resourceSchema) {
	return Object.keys(resourceSchema.properties).map((field) => {
		return schemaField(field, resourceSchema)
	}).join('\n')
}

function schemaField(name, resourceSchema) {
	const schema = resourceSchema.properties[name]
	switch (schema.type) {
		case 'string':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			`
		case 'number':
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
				Optional: true,
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
				Computed: true,
				Optional: true,
			},
			`
		case 'array':
			if (schema.items && schema.items.type === "object") {
				return `
				"${name}": &schema.Schema{
					Type: schema.TypeList,
					Computed: true,
					Optional: true,
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
					Computed: true,
					Optional: true,
				},
				`
			}
		case 'object':
		default:
			return `
			"${name}": &schema.Schema{
				Type: schema.TypeMap,
				Elem: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			`
	}
}
