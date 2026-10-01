const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const {spawnSync} = require("node:child_process");
const test = require("node:test");

test("Canvas compatibility normalization preserves #479 and remains idempotent", () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "rootly-clean-swagger-"));
  const swaggerPath = path.join(directory, "swagger.json");
  const channel = {
    type: "object",
    properties: {id: {type: "string"}, name: {type: "string"}},
    required: ["id", "name"],
  };
  const swagger = {
    paths: {},
    components: {schemas: {
      create_slack_canvas_task_params: {properties: {
        channel: structuredClone(channel),
        retry_count: {type: "integer"},
        retry_wait_time: {type: "integer"},
      }},
      update_slack_canvas_task_params: {properties: {
        channel: structuredClone(channel),
        operation: {type: "string", enum: ["insert_at_end", "replace"]},
        section_name: {type: "string"},
        retry_count: {type: "integer"},
        retry_wait_time: {type: "integer"},
      }},
    }},
  };
  try {
    fs.writeFileSync(swaggerPath, JSON.stringify(swagger));
    for (let run = 0; run < 2; run++) {
      const result = spawnSync(process.execPath, [path.join(__dirname, "clean-swagger.js"), swaggerPath], {encoding: "utf8"});
      assert.equal(result.status, 0, result.stderr);
    }
    const cleaned = JSON.parse(fs.readFileSync(swaggerPath, "utf8"));
    for (const action of ["create_slack_canvas", "update_slack_canvas"]) {
      const params = cleaned.components.schemas[`${action}_task_params`].properties;
      const annotated = params.channel.properties.workspace;
      assert.equal(params.channel.tf_nested_object, true);
      assert.equal(annotated.tf_nested_object, true);
      assert.equal(annotated.nullable, true);
      assert.equal(annotated.properties.id.tf_no_liquid, true);
      assert.equal((annotated.description.match(/Typed Go requests/g) || []).length, 1);
      assert.deepEqual([params.retry_count.minimum, params.retry_count.maximum], [0, 4]);
      assert.deepEqual([params.retry_wait_time.minimum, params.retry_wait_time.maximum], [1, 15]);
    }
    const update = cleaned.components.schemas.update_slack_canvas_task_params.properties;
    assert.deepEqual(update.operation.enum, ["insert_at_end", "replace", "managed_sections"]);
    assert.equal(update.section_name, undefined);
  } finally {
    fs.rmSync(directory, {recursive: true, force: true});
  }
});

test("status-page announcement fields match their operation contracts", () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "rootly-clean-swagger-"));
  const swaggerPath = path.join(directory, "swagger.json");
  const swagger = {
    paths: {},
    components: {schemas: {
      status_page_announcement: {properties: {
        user_id: {type: "integer"},
        published_at: {type: "string"},
      }},
      new_status_page_announcement: {properties: {data: {properties: {attributes: {properties: {
        notify_subscribers: {type: "boolean", description: "Notify subscribers"},
      }}}}}},
    }},
  };
  try {
    fs.writeFileSync(swaggerPath, JSON.stringify(swagger));
    const result = spawnSync(process.execPath, [path.join(__dirname, "clean-swagger.js"), swaggerPath], {encoding: "utf8"});
    assert.equal(result.status, 0, result.stderr);
    const properties = JSON.parse(fs.readFileSync(swaggerPath, "utf8"))
      .components.schemas.status_page_announcement.properties;
    assert.deepEqual(properties.notify_subscribers, {
      type: "boolean",
      description: "Notify subscribers",
      default: true,
      tf_create_only: true,
    });
    assert.equal(properties.user_id.tf_computed, true);
    assert.equal(properties.published_at.tf_computed, true);
  } finally {
    fs.rmSync(directory, {recursive: true, force: true});
  }
});

test("required workflow-task properties omitted by OpenAPI are restored", () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), "rootly-clean-swagger-"));
  const swaggerPath = path.join(directory, "swagger.json");
  const swagger = {
    paths: {},
    components: {schemas: {
      auto_assign_role_rootly_task_params: {
        type: "object",
        properties: {task_type: {type: "string"}},
        required: ["incident_role_id"],
      },
    }},
  };
  try {
    fs.writeFileSync(swaggerPath, JSON.stringify(swagger));
    const result = spawnSync(process.execPath, [path.join(__dirname, "clean-swagger.js"), swaggerPath], {encoding: "utf8"});
    assert.equal(result.status, 0, result.stderr);
    const property = JSON.parse(fs.readFileSync(swaggerPath, "utf8"))
      .components.schemas.auto_assign_role_rootly_task_params.properties.incident_role_id;
    assert.deepEqual(property, {type: "string", description: "The role id"});
  } finally {
    fs.rmSync(directory, {recursive: true, force: true});
  }
});
