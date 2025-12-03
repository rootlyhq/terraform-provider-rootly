import { camelize } from "inflection";
import { OpenAPIObject } from "./schema";

const RESOURCES = ["alert_route"];

async function getRootlySwagger() {
  const response = await fetch(
    "https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json"
  );
  return await response.json();
}

function generateResource(swagger: OpenAPIObject, resource: string) {
  const resourceSchema = swagger.components.schemas[resource];
  if (!resourceSchema) {
    throw new Error(`Resource ${resource} not found`);
  }

  const resourceName = `${camelize(resource)}Resource`;
  const modelName = `${camelize(resource)}Model`;

  return `
package provider

import (
  "context"
  "github.com/hashicorp/terraform-plugin-framework/resource"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema"
  "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*${resourceName})(nil)

func New${resourceName}() resource.Resource {
  return &${resourceName}{}
}

type ${resourceName} struct {
  baseResource
}

func (r *${resourceName}) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_${resource}"
}

func (r *${resourceName}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (r *${resourceName}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ${modelName}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ${modelName}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ${modelName}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ${modelName}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}

type ${modelName} struct {
	Id types.String
}
`;
}

async function main() {
  console.log("🚀 Fetching Rootly Swagger...");
  const originalSwagger = await getRootlySwagger();

  await Bun.write(
    new URL("swagger.json", import.meta.url),
    JSON.stringify(originalSwagger, null, 2)
  );

  const swagger = OpenAPIObject.parse(originalSwagger);

  await Bun.write(
    new URL("swagger.parsed.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  for (const resource of RESOURCES) {
    const code = generateResource(swagger, resource);
    await Bun.write(
      new URL(`../provider/resource_${resource}_gen.go`, import.meta.url),
      code
    );
  }
}

await main()
  .then(() => {
    console.log("✨ Done");
    process.exit(0);
  })
  .catch((e) => {
    console.error(e);
    process.exit(1);
  });
