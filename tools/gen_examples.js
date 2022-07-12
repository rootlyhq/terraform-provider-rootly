const fs = require("fs")
const path = require('path')
const templatesPath = process.argv[2]
let templates = null

if (!templatesPath) {
  console.error("No templates JSON specified. Use `node tools/gen_examples.js <templates_file_path>`")
  process.exit(1)
}

try {
  templates = require(path.resolve(templatesPath))
} catch (error) {
  throw new Error(`Error loading templates JSON: ${error.message}`)
}

templates.forEach((tpl) => {
  fs.writeFileSync(`./examples/provider/${tpl.template}.tf`, genExample(tpl))
})

function genExample(tpl) {
  const params = Object.keys(tpl.config).filter((key) => ["kind", "name", "description"].indexOf(key) === -1)
  return `resource "rootly_workflow_${tpl.config.kind}" "${tpl.template}" {
  name = "${tpl.config.name}"
  description = "${tpl.config.description}"
  trigger_params {
    ${params.map((key) => {
      const param = tpl.config[key]
      return genExampleTaskParam(key, param)
    }).join("\n    ")}
  }
  enabled = true
}

${tpl.tasks.map((task) => genExampleTask(tpl, task)).join("\n\n")}`
}

function genExampleTask(tpl, task) {
  const task_key = Object.keys(task)[0]
  const task_name = task_key.replace("genius_", "").replace("_task", "")
  task = task[task_key]
  return `resource "rootly_workflow_task_${task_name}" "${task_name}" {
  workflow_id = rootly_workflow_${tpl.config.kind}.${tpl.template}.id
  task_params {
    ${Object.keys(task).map((key) => genExampleTaskParam(key, task[key])).join("\n    ")}
  }
}`
}

function genExampleTaskParam(key, val) {
  if (typeof val === "string") {
    if (val.indexOf("\n") !== -1) {
      return `${key} = <<EOT\n${val}\nEOT`
    }
    return `${key} = "${val.replace(/"/g, '\\"')}"`
  }
  if (typeof val === "number") {
    return `${key} = ${val}`
  }
  if (Array.isArray(val)) {
    if (typeof val[0] === "object") {
      return val.map((v) => `${key} {
      id = "${v.id}"
      name = "${v.name}"
    }`).join("\n    ")
    }
    return `${key} = [${val.map((v) => `"${v}"`).join(", ")}]`
  }
  if (typeof val === "object") {
    return `${key} = {
      ${Object.keys(val).map((k) => genExampleTaskParam(k, val[k])).join("\n      ")}
    }`
  }
}
