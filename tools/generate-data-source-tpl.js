const inflect = require('./inflect');

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

func dataSource${nameCamel}() *schema.Resource {
	return &schema.Resource {
		ReadContext: dataSource${nameCamel}Read,
		Schema: map[string]*schema.Schema {
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
			},
			${schemaFields(resourceSchema, filterParameters)}
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

function filterCamelize(name) {
	return inflect.camelize(name.replace(/[\[\]]+/g, '_').replace(/_$/, ''))
}

function filterUnderscore(name) {
	return inflect.underscore(filterCamelize(name)).replace('filter_', '')
}

function setFilterFields(name, resourceSchema, filterParameters) {
	return (filterParameters || []).filter((paramSchema) => {
		return paramSchema.name.match(/^filter/)
	}).map((paramSchema) => {
		const filterField = filterUnderscore(paramSchema.name)
		const fieldSchema = resourceSchema.properties[filterField]
		if (fieldSchema) {
			return `
				if value, ok := d.GetOkExists("${filterField}"); ok {
					${filterField} := value.(${jsonapiToGoType(fieldSchema.type)})
					params.${filterCamelize(paramSchema.name)} = &${filterField}
				}
			`
		} else if (paramSchema.name.match(/(lt|gt)\]/)) {
			const rangeKey = filterField.split('_').pop()
			return `
				${filterField} := d.Get("${filterField.replace(/_(lt|gt)$/, '')}").(map[string]interface{})
				if value, exists := ${filterField}["${rangeKey}"]; exists {
					v := value.(string)
					params.${filterCamelize(paramSchema.name)} = &v
				}
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

function schemaFields(resourceSchema, filterParameters) {
	return Object.keys(resourceSchema.properties).filter((name) => {
		return filterParameters.some((param) => param.name.match(name))
	}).map((name) => {
		return schemaField(name, resourceSchema);
	}).join('\n')
}

function schemaField(name, resourceSchema, filterParameters) {
	const schema = resourceSchema.properties[name]
	switch (schema.type) {
		case 'integer':
		case 'number':
			return `
			"${name}": &schema.Schema {
				Type: schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			`
		case 'boolean':
			if (name === "enabled") {
				return `
				"${name}": &schema.Schema {
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
				},
				`
			}
			return `
			"${name}": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			`
		case 'string':
		default:
			if (name.match(/_at$/)) {
				return `
				"${name}": &schema.Schema {
					Type: schema.TypeMap,
					Description: "Filter by date range using 'lt' and 'gt'.",
					Optional: true,
				},
				`
			}
			return `
			"${name}": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Optional: true,
			},
			`
	}
}
