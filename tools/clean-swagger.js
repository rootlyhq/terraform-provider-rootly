const fs = require("fs");
const swagger = require(require("path").resolve(
  __dirname,
  "..",
  process.argv[2]
));

// Find filter paramters with incorrect "string" type and change to "bool"
function fixFilterParameterTypes(obj) {
  if (obj["/v1/incidents"]) {
    obj["/v1/incidents"].get.parameters.forEach((paramSchema) => {
      if (paramSchema.name === "filter[private]") {
        paramSchema.schema.type = 'boolean'
      }
    })
  }
}

// Terraform doesn't support polymorphism and oapi-codegen doesn't support anyOf
function stripAnyOf(obj) {
  if (typeof obj === "object" && obj !== null) {
    if (
      obj.hasOwnProperty("anyOf") &&
      obj.anyOf[0].required
    ) {
      delete obj.anyOf;
    }
    Object.keys(obj).forEach(function (key) {
      stripAnyOf(obj[key]);
    });
  }
}

// Terraform doesn't support polymorphism and oapi-codegen doesn't support oneOf
// So we look for oneOf arrays of objects and merge them into a single object
function combineOneOf(obj) {
  if (typeof obj === "object" && obj !== null) {
    if (obj.hasOwnProperty("oneOf") && obj.oneOf[0] && obj.oneOf[0].properties) {
      obj.properties = obj.oneOf.reduce((accum, oneOfItem) => {
        if (oneOfItem.properties) {
          Object.keys(oneOfItem.properties).forEach((key) => {
            accum[key] = oneOfItem.properties[key]
          })
        }
        return accum
      }, {})
      delete obj.oneOf
    }
    Object.keys(obj).forEach(function (key) {
      combineOneOf(obj[key]);
    });
  }
}

function renameEscalationPolicyLevelSchemas(obj) {
  for (var key in obj) {
    if (key.match(/escalation_policy_level/)) {
      let newKey = key.replace(/escalation_policy_level/g, "escalation_level")
      obj[newKey] = obj[key]
      delete obj[key];
      key = newKey
    }
    if (typeof obj[key] === "string" && obj[key].match(/components\/schemas/)) {
      obj[key] = obj[key].replace(/escalation_policy_level/g, "escalation_level")
    } else if (typeof obj[key] === "object" && obj[key] !== null) {
      renameEscalationPolicyLevelSchemas(obj[key])
    }
  }
}

fixFilterParameterTypes(swagger.paths);
stripAnyOf(swagger.components.schemas);
combineOneOf(swagger.components.schemas);
renameEscalationPolicyLevelSchemas(swagger);
fs.writeFileSync(process.argv[2], JSON.stringify(swagger));
