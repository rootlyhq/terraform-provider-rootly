const fs = require("fs")
const path = require('path')
const swaggerPath = process.argv[2]
const inflect = require('inflect')
const providerTpl = require('./generate-provider-tpl')
const clientTpl = require('./generate-client-tpl')
const dataSourceTpl = require('./generate-data-source-tpl')
const resourceTpl = require('./generate-resource-tpl')
const generateTasks = require('./generate-tasks')
const swagger = require(path.resolve(swaggerPath))

const excluded = [
	"dashboard",
	"dashboard_panel",
	"workflow",
	"workflow_task",
	"incident",
	"incident_post_mortem",
	"incident_action_item",
	"incident_event",
	"playbook",
	"playbook_task",
	"incident_role_task",
	"incident_feedback",
	"incident_custom_field_selection",
	"workflow_custom_field_selection",
	"status_page_template",
	"pulse",
	"alert",
]

const resources = Object.keys(swagger.components.schemas).filter((name) => {
	return excluded.indexOf(name) === -1 && collectionPathSchema(name)
})

resources.forEach(generate)

generateTasks(swagger)
generateProvider(resources)

function generate(name) {
	generateClient(name)
	generateDataSource(name)
	generateResource(name)
	generateResourceTest(name)
}

function generateProvider(resources) {
	const code = providerTpl(resources)
	fs.writeFileSync(path.resolve(__dirname, '..', 'provider', 'provider.go'), code)
}

function generateClient(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	const code = clientTpl(name, resourceSchema(name), pathIdField)
	fs.writeFileSync(path.resolve(__dirname, '..', 'client', `${inflect.pluralize(name)}.go`), code)
}

function generateDataSource(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	const code = dataSourceTpl(name, resourceSchema(name), collectionSchema, pathIdField)
	fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `data_source_${name}.go`), code)
}

function generateResource(name) {
	const collectionSchema = collectionPathSchema(name)
	const pathIdField = collectionSchema && collectionSchema.parameters && collectionSchema.parameters[0] && collectionSchema.parameters[0].name
	const code = resourceTpl(name, resourceSchema(name), collectionSchema, pathIdField, createResourceSchema(name))
	fs.writeFileSync(path.resolve(__dirname, '..', 'provider', `resource_${name}.go`), code)
}

function generateResourceTest(name) {
}

function resourceSchema(name) {
	return swagger.components.schemas[name]
}

function createResourceSchema(name) {
	return swagger.components.schemas[`new_${name}`]
}

function collectionPathSchema(name) {
	return Object.keys(swagger.paths).filter((url) => {
		const get = swagger.paths[url].get
		return get && get.operationId.replace(/ /g, '') === `list${inflect.pluralize(inflect.camelize(name))}`
	}).map((url) => swagger.paths[url])[0]
}
