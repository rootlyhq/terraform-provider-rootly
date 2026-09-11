const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");
const generateTasks = require("./generate-tasks");

test("nested objects opt in without changing flat workflow task maps", () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "rootly-task-generator-"));
  const previous = process.cwd();
  try {
    fs.mkdirSync(path.join(directory, "provider"));
    process.chdir(directory);
    const workspace = {type: "object", tf_nested_object: true, description: "API description", tf_description: "Terraform description", properties: {id: {type: "string"}, name: {type: "string"}}, required: ["id", "name"]};
    const channel = {type: "object", properties: {id: {type: "string"}, name: {type: "string"}}};
    generateTasks(["flat", "nested"], {components: {schemas: {
      flat_task_params: {properties: {channel}, required: ["channel"]},
      nested_task_params: {properties: {retry_count: {type: "integer", default: 0}, retry_wait_time: {type: "integer", default: 1}, channel: {...channel, tf_nested_object: true, properties: {...channel.properties, workspace}}}, required: ["channel"]},
    }}});
    const flat = fs.readFileSync("provider/resource_workflow_task_flat.go", "utf8");
    const nested = fs.readFileSync("provider/resource_workflow_task_nested.go", "utf8");
    assert.match(flat, /Type: schema.TypeMap/);
    assert.doesNotMatch(flat, /FlattenWorkflowTaskObjects|ExpandWorkflowTaskObjects/);
    assert.match(nested, /"workspace": &schema.Schema[\s\S]*?Type: schema.TypeList/);
    assert.equal((nested.match(/FlattenWorkflowTaskObjects/g) || []).length, 2);
    assert.equal((nested.match(/ExpandWorkflowTaskObjects/g) || []).length, 1);
    assert.doesNotMatch(nested, /Map must contain two fields/);
    assert.match(nested, /_, err = c.UpdateWorkflowTask/);
    assert.match(nested, /Terraform description/);
    assert.match(nested, /"retry_count": &schema.Schema[\s\S]*?Default: 0/);
    assert.match(nested, /"retry_wait_time": &schema.Schema[\s\S]*?Default: 1/);
    assert.doesNotMatch(nested, /API description/);
    assert.match(nested, /if err := d.Set\("task_params"/);
  } finally {
    process.chdir(previous);
    fs.rmSync(directory, {recursive: true, force: true});
  }
});
