const assert = require("node:assert/strict");
const fs = require("node:fs");
const {spawnSync} = require("node:child_process");
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
    const workspace = {type: "object", tf_nested_object: true, description: "API description", tf_description: "Terraform description", properties: {id: {type: "string", minLength: 1, pattern: "\\S"}, name: {type: "string", minLength: 1, pattern: "\\S"}}, required: ["id", "name"]};
    const channel = {type: "object", properties: {id: {type: "string", minLength: 1, pattern: "\\S"}, name: {type: "string", minLength: 1, pattern: "\\S"}}, required: ["id", "name"]};
    generateTasks(["flat", "nested"], {components: {schemas: {
      flat_task_params: {properties: {channel}, required: ["channel"]},
      nested_task_params: {properties: {retry_count: {type: "integer", default: 0, minimum: 0, maximum: 4}, retry_wait_time: {type: "integer", default: 1, minimum: 1, maximum: 15}, channel: {...channel, tf_nested_object: true, properties: {...channel.properties, workspace}}}, required: ["channel"]},
    }}});
    const flat = fs.readFileSync("provider/resource_workflow_task_flat.go", "utf8");
    const nested = fs.readFileSync("provider/resource_workflow_task_nested.go", "utf8");
    const flatTests = fs.readFileSync("provider/resource_workflow_task_flat_test.go", "utf8");
    const nestedTests = fs.readFileSync("provider/resource_workflow_task_nested_test.go", "utf8");
    assert.match(flat, /Type: schema.TypeMap/);
    assert.doesNotMatch(flat, /FlattenWorkflowTaskObjects|ExpandWorkflowTaskObjects/);
    assert.match(nested, /"workspace": &schema.Schema[\s\S]*?Type: schema.TypeList/);
    assert.equal((nested.match(/FlattenWorkflowTaskObjects/g) || []).length, 2);
    assert.equal((nested.match(/ExpandWorkflowTaskObjects/g) || []).length, 1);
    assert.doesNotMatch(nested, /Map must contain two fields/);
    assert.match(nested, /_, err = c.UpdateWorkflowTask/);
    assert.match(nested, /Terraform description/);
    assert.equal((nested.match(/ValidateFunc: validation.StringIsNotWhiteSpace/g) || []).length, 4);
    assert.doesNotMatch(flat, /validation.StringIsNot(?:Empty|WhiteSpace)/);
    assert.match(nested, /"retry_count": &schema.Schema[\s\S]*?Default: 0/);
    assert.match(nested, /"retry_wait_time": &schema.Schema[\s\S]*?Default: 1/);
    assert.match(nested, /"retry_count": &schema.Schema[\s\S]*?ValidateFunc: validation.IntBetween\(0, 4\)/);
    assert.match(nested, /"retry_wait_time": &schema.Schema[\s\S]*?ValidateFunc: validation.IntBetween\(1, 15\)/);
    for (const [source, isNested] of [[flatTests, false], [nestedTests, true]]) {
      const configs = [...source.matchAll(/return fmt\.Sprintf\(`([\s\S]*?)`, name\)/g)];
      assert.equal(configs.length, 2, "create and update test configs must be generated");
      for (const [, template] of configs) {
        const config = template.replaceAll("%s", "test-workflow").replaceAll("%%", "%");
        if (isNested) {
          assert.match(config, /channel\s*\{\s*id\s*=\s*"test"\s*name\s*=\s*"test"\s*workspace\s*\{\s*id\s*=\s*"test"\s*name\s*=\s*"test"\s*\}\s*\}/);
          assert.doesNotMatch(config, /\b(?:channel|workspace)\s*=/);
        } else {
          assert.match(config, /channel\s*=\s*\{\s*id\s*=\s*"foo"\s*name\s*=\s*"bar"\s*\}/);
        }
        const formatted = spawnSync("terraform", ["fmt", "-write=false", "-no-color", "-"], {input: config, encoding: "utf8"});
        assert.ifError(formatted.error);
        assert.equal(formatted.status, 0, formatted.stderr);
      }
    }
    assert.doesNotMatch(nested, /API description/);
    assert.match(nested, /if err := d.Set\("task_params"/);
  } finally {
    process.chdir(previous);
    fs.rmSync(directory, {recursive: true, force: true});
  }
});


test("nested JSON fields generate compilable imports without changing flat maps", () => {
  const previous = process.cwd();
  const directory = fs.mkdtempSync(path.join(previous, ".workflow-codegen-"));
  try {
    fs.mkdirSync(path.join(directory, "provider"));
    process.chdir(directory);
    const channel = {type: "object", properties: {
      id: {type: "string"},
      name: {type: "string"},
      settings: {type: "object", tf_nested_object: true, properties: {
        payload: {type: "string", description: "JSON payload"},
      }},
    }, required: ["id", "name"]};
    generateTasks(["nested_json", "flat_json"], {components: {schemas: {
      nested_json_task_params: {properties: {channel: {...channel, tf_nested_object: true}}, required: ["channel"]},
      flat_json_task_params: {properties: {channel}, required: ["channel"]},
    }}});
    const nested = fs.readFileSync("provider/resource_workflow_task_nested_json.go", "utf8");
    const flat = fs.readFileSync("provider/resource_workflow_task_flat_json.go", "utf8");
    assert.match(nested, /"encoding\/json"/);
    assert.match(nested, /"reflect"/);
    assert.match(nested, /json\.Unmarshal/);
    assert.match(nested, /reflect\.DeepEqual/);
    assert.doesNotMatch(flat, /"encoding\/json"|"reflect"|json\.Unmarshal|reflect\.DeepEqual/);
    for (const task of ["nested_json", "flat_json"]) {
      fs.unlinkSync(`provider/resource_workflow_task_${task}_test.go`);
    }
    fs.writeFileSync("provider/validation.go", `package provider
import (
  "context"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
func validateUniqueWorkflowTaskPosition(context.Context, *schema.ResourceDiff, interface{}) error { return nil }
`);
    const compiled = spawnSync("go", ["test", `./${path.basename(directory)}/provider`], {cwd: previous, encoding: "utf8"});
    assert.ifError(compiled.error);
    assert.equal(compiled.status, 0, compiled.stdout + compiled.stderr);
  } finally {
    process.chdir(previous);
    fs.rmSync(directory, {recursive: true, force: true});
  }
});
