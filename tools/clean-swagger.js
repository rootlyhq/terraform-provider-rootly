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
      obj.anyOf[0].required &&
      !obj.properties
    ) {
      obj.properties = {};
      obj.anyOf.forEach((child) => {
        if (child.properties) {
          Object.keys(child.properties).forEach((child_property_key) => {
            if (obj.properties[child_property_key]) {
              if (obj.properties[child_property_key].type === "string" && obj.properties[child_property_key].enum) {
                obj.properties[child_property_key].enum = obj.properties[child_property_key].enum.concat(child.properties[child_property_key].enum)
              }
            } else {
              obj.properties[child_property_key] = child.properties[child_property_key]
            }
            obj.properties[child_property_key].anyOfChild = true
          })
        }
      })
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

function renameEscalationPolicyPathSchemas(obj) {
  for (var key in obj) {
    let value = obj[key];

    // TODO: Ignore escalation_policy_(path|level) in codegen and define manually
    // All the naming for the escalation API resources is inconsistent
    if (key.match(/escalation_policy_path/) && !obj.delay) {
      let newKey = key.replace(/escalation_policy_path/g, "escalation_path")
      obj[newKey] = obj[key]
      delete obj[key];
      key = newKey
    }

    if (typeof obj[key] === "string" && obj[key].match(/components\/schemas/)) {
      obj[key] = obj[key].replace(/escalation_policy_path/g, "escalation_path")
    } else if (typeof obj[key] === "object" && obj[key] !== null) {
      renameEscalationPolicyPathSchemas(obj[key])
    }

    if (typeof value === "string" && value.match(/escalation_policy_path/)) {
      obj[key] = value.replace(/escalation_policy_path/, "escalation_path")
    }
  }
}

function stripParameterAnyOf(obj) {
  if (typeof obj === "object" && obj !== null) {
    if (obj.anyOf && obj.anyOf.every(item => item.type === "string")) {
      return {"type": "string"};
    }
    return obj;
  }
}

fixFilterParameterTypes(swagger.paths);
stripAnyOf(swagger.components.schemas);
combineOneOf(swagger.components.schemas);
stripParameterAnyOf(swagger.components.parameters);
renameEscalationPolicyLevelSchemas(swagger);
renameEscalationPolicyPathSchemas(swagger);
fs.writeFileSync(process.argv[2], JSON.stringify(swagger));
