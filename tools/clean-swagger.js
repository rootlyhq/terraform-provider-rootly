const fs = require('fs')
const swagger = require(require('path').resolve(__dirname, '..', process.argv[2]))

// strip anyOf from schema, oapi-codegen fails if it encounters it
function strip(obj) {
	if (typeof obj === "object" && obj !== null) {
		if (obj.hasOwnProperty("anyOf") && obj.anyOf[0].required && Object.keys(obj.anyOf[0]).length === 1) {
			delete obj.anyOf
		}
		Object.keys(obj).forEach(function(key) {
			strip(obj[key])
		})
	}
}

strip(swagger.components.schemas)
fs.writeFileSync(process.argv[2], JSON.stringify(swagger))
