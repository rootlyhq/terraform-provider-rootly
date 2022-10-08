const inflect = require('inflect')

module.exports = {
	...inflect,
	camelize: (name) => {
		if (name === "post_mortem_template") {
			return "PostmortemTemplate"
		}
		return inflect.camelize(name)
	},
}
