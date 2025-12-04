import {
  camelize,
  humanize,
  pluralize,
  singularize,
  underscore,
} from "inflection";
import { match } from "ts-pattern";
import { toIR, type IRObject, type IRResource, type IRType } from "./ir";
import { SWAGGER_MODS } from "./swagger-mods";

const RESOURCES = ["alert_route", "dashboard_panel"];

async function getRootlySwagger() {
  const response = await fetch(
    "https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json"
  );
  return (await response.json()) as any;
}

function generateTerraformValuer({
  prefix,
  parent,
  field,
  ir,
}: {
  prefix: string;
  parent: string;
  field: string;
  ir: IRType;
}): {
  output: string;
  hasErr: boolean;
} {
  return match(ir)
    .returnType<{ output: string; hasErr: boolean }>()
    .with({ kind: "string" }, () => ({
      output: `${prefix}${camelize(field)}.ValueString()`,
      hasErr: false,
    }))
    .with({ kind: "bool" }, () => ({
      output: `${prefix}${camelize(field)}.ValueBool()`,
      hasErr: false,
    }))
    .with({ kind: "int" }, () => ({
      output: `${prefix}${camelize(field)}.ValueInt64()`,
      hasErr: false,
    }))
    .with({ kind: "array", element: { kind: "object" } }, () => ({
      output: `
        func() ([]apiclient.${camelize(`${parent}_${field}_item`)}, error) {
          var itemClientModels []apiclient.${camelize(
            `${parent}_${field}_item`
          )}
          for _, item := range ${prefix}${camelize(field)}.MustGet(ctx) {
            itemClientModel, err := item.ToClientModel(ctx)
            if err != nil {
              return nil, err
            }
            itemClientModels = append(itemClientModels, *itemClientModel)
          }
          return itemClientModels, nil
        }()
      `.trim(),
      hasErr: true,
    }))
    .with({ kind: "array" }, () => ({
      output: `${prefix}${camelize(field)}.MustGet(ctx)`,
      hasErr: false,
    }))
    .with({ kind: "object" }, () => ({
      output: `
        func() (apiclient.${camelize(`${parent}_${field}`)}, error) {
          model, diags := ${prefix}${camelize(field)}.Get(ctx)
          if diags.HasError() {
            return apiclient.${camelize(
              `${parent}_${field}`
            )}{}, fmt.Errorf("%v", diags.Errors())
          } else if model == nil {
            return apiclient.${camelize(
              `${parent}_${field}`
            )}{}, fmt.Errorf("model is nil")
          }
          clientModel, err := model.ToClientModel(ctx)
          if err != nil {
            return apiclient.${camelize(`${parent}_${field}`)}{}, err
          }
          return *clientModel, nil
        }()
      `.trim(),
      hasErr: true,
    }))
    .otherwise(() => {
      throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
    });
}

function generateTerraformType({
  parent,
  field,
  ir,
  inObject,
}: {
  parent: string;
  field: string;
  ir: IRType;
  inObject?: boolean;
}): { output: string; nested: string[] } {
  return match(ir)
    .returnType<{ output: string; nested: string[] }>()
    .with({ kind: "string" }, () => ({
      output: "supertypes.StringValue",
      nested: [],
    }))
    .with({ kind: "bool" }, () => ({
      output: "supertypes.BoolValue",
      nested: [],
    }))
    .with({ kind: "int" }, () => ({
      output: "supertypes.Int64Value",
      nested: [],
    }))
    .with({ kind: "array", element: { kind: "object" } }, (ir) => {
      const inner = generateTerraformType({
        inObject: true,
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `supertypes.ListNestedObjectValueOf[${inner.output}]`,
        nested: inner.nested,
      };
    })
    .with({ kind: "array" }, (ir) => {
      const inner = generatePrimitiveType({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `supertypes.ListValueOf[${inner.output}]`,
        nested: inner.nested,
      };
    })
    .with({ kind: "object" }, (ir) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateModel({ name: structName, ir });
      return {
        output: inObject
          ? structName
          : `supertypes.SingleNestedObjectValueOf[${structName}]`,
        nested: [struct],
      };
    })
    .otherwise(() => {
      throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
    });
}

function generatePrimitiveType({
  parent,
  field,
  ir,
}: {
  parent: string;
  field: string;
  ir: IRType;
}): { output: string; nested: string[] } {
  return match(ir)
    .returnType<{ output: string; nested: string[] }>()
    .with({ kind: "string" }, () => ({
      output: "string",
      nested: [],
    }))
    .with({ kind: "bool" }, () => ({
      output: "bool",
      nested: [],
    }))
    .with({ kind: "int" }, () => ({
      output: "int64",
      nested: [],
    }))
    .with({ kind: "array", element: { kind: "object" } }, (ir) => {
      const inner = generatePrimitiveType({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `[]${inner.output}`,
        nested: inner.nested,
      };
    })
    .with({ kind: "array" }, (ir) => {
      const inner = generatePrimitiveType({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `[]${inner.output}`,
        nested: inner.nested,
      };
    })
    .with({ kind: "object" }, (ir) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateClientModel({ name: structName, ir });
      return {
        output: structName,
        nested: [struct],
      };
    })
    .otherwise(() => {
      throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
    });
}

function generateModelValuer({
  parent,
  field,
  ir,
}: {
  parent: string;
  field: string;
  ir: IRType;
}): {
  output: string;
  nested: string[];
  hasErr: boolean;
} {
  return match(ir)
    .returnType<{
      output: string;
      nested: string[];
      hasErr: boolean;
    }>()
    .with({ kind: "string" }, () => ({
      output: `supertypes.NewStringValueOrNull(in.${camelize(field)})`,
      nested: [],
      hasErr: false,
    }))
    .with({ kind: "bool" }, () => ({
      output: `supertypes.NewBoolValue(in.${camelize(field)})`,
      nested: [],
      hasErr: false,
    }))
    .with({ kind: "int" }, () => ({
      output: `supertypes.NewInt64Value(in.${camelize(field)})`,
      nested: [],
      hasErr: false,
    }))
    .with({ kind: "array", element: { kind: "object" } }, (ir) => {
      const inner = generateModelValuer({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      const itemType = camelize(`${parent}_${field}_item`);
      return {
        output: `
          func() (supertypes.ListNestedObjectValueOf[${itemType}], error) {
            var elements []${itemType}
            for _, item := range in.${camelize(field)} {
              var element ${itemType}
              err := Fill${itemType}(ctx, item, &element)
              if err != nil {
                return supertypes.NewListNestedObjectValueOfNull[${itemType}](ctx), err
              }
              elements = append(elements, element)
            }
            return supertypes.NewListNestedObjectValueOfValueSlice(ctx, elements), nil
          }()
        `.trim(),
        nested: inner.nested,
        hasErr: true,
      };
    })
    .with({ kind: "array" }, (ir) => {
      const inner = generateModelValuer({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `supertypes.NewListValueOfSlice(ctx, in.${camelize(field)})`,
        nested: inner.nested,
        hasErr: false,
      };
    })
    .with({ kind: "object" }, (ir) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateFillModel({ name: structName, ir });
      return {
        output: `
          func() (supertypes.SingleNestedObjectValueOf[${structName}], error) {
            var out ${structName}
            err := Fill${structName}(ctx, in.${camelize(field)}, &out)
            if err != nil {
              return supertypes.NewSingleNestedObjectValueOfNull[${structName}](ctx), err
            }
            return supertypes.NewSingleNestedObjectValueOf(ctx, &out), nil
          }()
        `.trim(),
        nested: [struct],
        hasErr: true,
      };
    })
    .otherwise(() => {
      throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
    });
}

function generateModel({
  name,
  ir,
}: {
  name: string;
  ir: IRResource | IRObject;
}) {
  const modelStructLines: string[] = [];
  const toClientModelLines: string[] = [];
  const nested: string[] = [];

  if (ir.kind === "resource") {
    const terraformType = generateTerraformType({
      parent: name,
      field: "id",
      ir: ir.idElement,
    });
    modelStructLines.push(`Id ${terraformType.output} \`tfsdk:"id"\``);
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const terraformType = generateTerraformType({
      parent: name,
      field,
      ir: fieldIR,
    });
    const terraformValuer = generateTerraformValuer({
      prefix: "m.",
      parent: name,
      field,
      ir: fieldIR,
    });

    nested.push(...terraformType.nested);
    modelStructLines.push(
      `${camelize(field)} ${terraformType.output} \`tfsdk:"${underscore(
        field
      )}"\``
    );

    toClientModelLines.push(`if !m.${camelize(field)}.IsNull() {`);
    if (terraformValuer.hasErr) {
      toClientModelLines.push(
        `
          var err error
          out.${camelize(field)}, err = ${terraformValuer.output}
          if err != nil {
            return nil, err
          }
        `.trim()
      );
    } else {
      toClientModelLines.push(
        `out.${camelize(field)} = ${terraformValuer.output}`
      );
    }
    toClientModelLines.push(`}\n`);
  }

  const struct = `
type ${name} struct {
  ${modelStructLines.join("\n")}
}

func (m *${name}) ToClientModel(ctx context.Context) (*apiclient.${name}, error) {
  var out apiclient.${name}
  ${toClientModelLines.join("\n")}
  return &out, nil
}
`;

  return [struct, ...nested].join("\n");
}

function generateFillModel({ name, ir }: { name: string; ir: IRType }) {
  if (ir.kind !== "object" && ir.kind !== "resource") {
    throw new Error("Model root must be an object");
  }

  const fillModelLines: string[] = [];
  const nested: string[] = [];

  if (ir.kind === "resource") {
    const modelValuer = generateModelValuer({
      parent: name,
      field: "id",
      ir: ir.idElement,
    });
    fillModelLines.push(`out.Id = ${modelValuer.output}`);
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const modelValuer = generateModelValuer({
      parent: name,
      field,
      ir: fieldIR,
    });

    nested.push(...modelValuer.nested);

    if (modelValuer.hasErr) {
      fillModelLines.push("{");
      fillModelLines.push("var err error");
      fillModelLines.push(
        `out.${camelize(field)}, err = ${modelValuer.output}`
      );
      fillModelLines.push("if err != nil {");
      fillModelLines.push("return err");
      fillModelLines.push("}");
      fillModelLines.push("}");
    } else {
      fillModelLines.push(`out.${camelize(field)} = ${modelValuer.output}`);
    }
  }

  const struct = `
func Fill${name}(ctx context.Context, in apiclient.${name}, out *${name}) error {
  ${fillModelLines.join("\n")}
  return nil
}
`;

  return [struct, ...nested].join("\n");
}

function generateAttribute({
  parent,
  field,
  ir,
}: {
  parent: string;
  field: string;
  ir: Exclude<IRType, IRResource>;
}): string {
  const commonLines: string[] = [];
  if (ir.computedOptionalRequired === "required") {
    commonLines.push("Required: true,");
  }
  if (ir.computedOptionalRequired === "optional") {
    commonLines.push("Optional: true,");
  }
  if (ir.computedOptionalRequired === "computed") {
    commonLines.push("Computed: true,");
  }
  if (ir.description) {
    commonLines.push(`Description: ${JSON.stringify(ir.description)},`);
  }
  if (ir.description) {
    commonLines.push(`MarkdownDescription: ${JSON.stringify(ir.description)},`);
  }
  const common = commonLines.join("\n");

  return match(ir)
    .returnType<string>()
    .with({ kind: "string" }, () => {
      return `schema.StringAttribute{
          ${common}
          CustomType: supertypes.StringType{},
        }`;
    })
    .with({ kind: "bool" }, () => {
      return `schema.BoolAttribute{
          ${common}
          CustomType: supertypes.BoolType{},
        }`;
    })
    .with({ kind: "int" }, () => {
      return `schema.Int64Attribute{
          ${common}
          CustomType: supertypes.Int64Type{},
        }`;
    })
    .with({ kind: "array", element: { kind: "object" } }, (ir) => {
      const structType = camelize(`${parent}_${field}_item`);

      const attrs = Object.entries(ir.element.fields)
        .map(([fieldName, fieldIR]) => {
          return `"${underscore(fieldName)}": ${generateAttribute({
            parent: structType,
            field: fieldName,
            ir: fieldIR,
          })},`;
        })
        .join("\n");

      return `schema.ListNestedAttribute{
          ${common}
          CustomType: supertypes.NewListNestedObjectTypeOf[${structType}](ctx),
          NestedObject: schema.NestedAttributeObject{
            Attributes: map[string]schema.Attribute{
              ${attrs}
            },
          },
        }`;
    })
    .with({ kind: "array" }, ({ element }) => {
      const primitiveType = generatePrimitiveType({
        parent,
        field,
        ir: element,
      });

      return `schema.ListAttribute{
          ${common}
          CustomType: supertypes.NewListTypeOf[${primitiveType.output}](ctx),
        }`;
    })
    .with({ kind: "object" }, ({ fields }) => {
      const structType = camelize(`${parent}_${field}`);

      const attrs = Object.entries(fields)
        .map(([fieldName, fieldIR]) => {
          return `"${underscore(fieldName)}": ${generateAttribute({
            parent: structType,
            field: fieldName,
            ir: fieldIR,
          })},`;
        })
        .join("\n");

      return `schema.SingleNestedAttribute{
          ${common}
          CustomType: supertypes.NewSingleNestedObjectTypeOf[${structType}](ctx),
          Attributes: map[string]schema.Attribute{
            ${attrs}
          },
        }`;
    })
    .otherwise(() => {
      throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
    });
}

function generateResourceSchema({ name, ir }: { name: string; ir: IRType }) {
  if (ir.kind !== "resource") {
    throw new Error("Resource schema must be a resource");
  }

  const attrs: string[] = [];

  attrs.push(
    `"id": ${generateAttribute({
      parent: name,
      field: "id",
      ir: ir.idElement,
    })},`
  );

  attrs.push(
    ...Object.entries(ir.fields).map(
      ([fieldName, fieldIR]) =>
        `"${underscore(fieldName)}": ${generateAttribute({
          parent: name,
          field: fieldName,
          ir: fieldIR,
        })},`
    )
  );

  return `
schema.Schema{
  Attributes: map[string]schema.Attribute{
    ${attrs.join("\n")}
  },
}`;
}

function generateResourceUpdateDiff({
  ir,
  name,
}: {
  ir: IRResource;
  name: string;
}) {
  const lines: string[] = [];

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const terraformValuer = generateTerraformValuer({
      prefix: "plan.",
      parent: name,
      field,
      ir: fieldIR,
    });
    lines.push(`if !data.${camelize(field)}.Equal(plan.${camelize(field)}) {`);

    if (terraformValuer.hasErr) {
      lines.push(
        `
          var err error
          modelIn.${camelize(field)}, err = ${terraformValuer.output}
          if err != nil {
            resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to convert plan to model: %v", err))
            return
          }
        `.trim()
      );
    } else {
      lines.push(`modelIn.${camelize(field)} = ${terraformValuer.output}`);
    }

    lines.push(`}\n`);
  }

  return lines.join("\n");
}

function generateResource({ ir, name }: { ir: IRResource; name: string }) {
  const baseName = camelize(name);
  const resourceName = `${baseName}Resource`;
  const modelName = `${baseName}Model`;

  return `
// DO NOT MODIFY: This file is generated by internal/providergen/index.ts. Any changes will be overwritten during the next build.
package provider

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/DataDog/jsonapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/apiclient"
)

var _ resource.Resource = (*${resourceName})(nil)
var _ resource.ResourceWithImportState = (*${resourceName})(nil)

func New${resourceName}() resource.Resource {
  return &${resourceName}{}
}

type ${resourceName} struct {
  baseResource
}

func (r *${resourceName}) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_${name}"
}

func (r *${resourceName}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = ${generateResourceSchema({ name: modelName, ir })}
}

func (r *${resourceName}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  var data ${modelName}

  // Read Terraform plan data into the model
  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

  if resp.Diagnostics.HasError() {
    return
  }

  // Create API call logic
  modelIn, err := data.ToClientModel(ctx)
  if err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to client model: %v", err))
    return
  }

  b, err := jsonapi.Marshal(&modelIn, jsonapi.MarshalClientMode())
  if err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to marshal model to JSON: %v", err))
    return
  }

  httpResp, err := r.client.Create${baseName}WithBodyWithResponse(
    ctx,
    ${
      ir.listPathIdParam
        ? `${
            generateTerraformValuer({
              prefix: "data.",
              parent: "",
              field: ir.listPathIdParam.name,
              ir: ir.listPathIdParam.element,
            }).output
          },`
        : ""
    }
    "application/vnd.api+json",
    bytes.NewReader(b),
  )
  if err != nil {
    resp.Diagnostics.AddError("API Error", err.Error())
    return
  } else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
    resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create, got status code: %d", httpResp.StatusCode()))
    return
  } else if httpResp.Body == nil {
    resp.Diagnostics.AddError("API Error", "Unable to create, got empty response")
    return
  }

  var modelOut apiclient.${modelName}
  if err := jsonapi.Unmarshal(httpResp.Body, &modelOut); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to unmarshal response: %v", err))
    return
  }

  if err := Fill${modelName}(ctx, modelOut, &data); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to fill model: %v", err))
    return
  }

  if resp.Diagnostics.HasError() {
    return
  }

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
  httpResp, err := r.client.Get${baseName}WithResponse(
    ctx,
    data.Id.ValueString(),
    ${ir.getHasQueryParams ? "nil," : ""}
  )
  if err != nil {
    resp.Diagnostics.AddError("API Error", err.Error())
    return
  } else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
    resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code: %d", httpResp.StatusCode()))
    return
  } else if httpResp.Body == nil {
    resp.Diagnostics.AddError("API Error", "Unable to read, got empty response")
    return
  }

  var modelOut apiclient.${modelName}
  if err := jsonapi.Unmarshal(httpResp.Body, &modelOut); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to unmarshal response: %v", err))
    return
  }

  if err := Fill${modelName}(ctx, modelOut, &data); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to fill model: %v", err))
    return
  }

  if resp.Diagnostics.HasError() {
    return
  }

  // Save updated data into Terraform state
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *${resourceName}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  var plan, data ${modelName}

  // Read Terraform plan data into the model
  resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

  if resp.Diagnostics.HasError() {
    return
  }

  // Read Terraform state data into the model
  resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

  if resp.Diagnostics.HasError() {
    return
  }

  // Update API call logic
  var modelIn apiclient.${modelName}

  ${generateResourceUpdateDiff({ name: modelName, ir })}

  b, err := jsonapi.Marshal(&modelIn, jsonapi.MarshalClientMode())
  if err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to marshal model to JSON: %v", err))
    return
  }

	httpResp, err := r.client.Update${baseName}WithBodyWithResponse(
    ctx,
    data.Id.ValueString(),
    "application/vnd.api+json",
    bytes.NewReader(b),
  )
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	} else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update, got status code: %d", httpResp.StatusCode()))
		return
  } else if httpResp.Body == nil {
    resp.Diagnostics.AddError("API Error", "Unable to read, got empty response")
    return
  }

  var modelOut apiclient.${modelName}
  if err := jsonapi.Unmarshal(httpResp.Body, &modelOut); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to unmarshal response: %v", err))
    return
  }

  if err := Fill${modelName}(ctx, modelOut, &data); err != nil {
    resp.Diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to fill model: %v", err))
    return
  }

  if resp.Diagnostics.HasError() {
    return
  }

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
  httpResp, err := r.client.Delete${baseName}WithResponse(
    ctx,
    data.Id.ValueString(),
  )
  if err != nil {
    resp.Diagnostics.AddError("API Error", err.Error())
    return
  } else if httpResp.StatusCode() == http.StatusNotFound {
    return
  } else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
    resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete, got status code: %d", httpResp.StatusCode()))
    return
  }
}

func (r *${resourceName}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
  resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

${generateModel({ name: modelName, ir })}
${generateFillModel({ name: modelName, ir })}
`;
}

function generateClientModel({ ir, name }: { ir: IRType; name: string }) {
  if (ir.kind !== "resource" && ir.kind !== "object") {
    throw new Error("Model root must be a resource");
  }

  const modelStructLines: string[] = [];
  const fillModelLines: string[] = [];
  const nested: string[] = [];

  const modelName =
    ir.kind === "resource" ? `${camelize(name)}Model` : camelize(name);

  if (ir.kind === "resource") {
    const primitiveType = generatePrimitiveType({
      parent: modelName,
      field: "ir",
      ir: ir.idElement,
    });
    modelStructLines.push(
      `Id ${primitiveType.output} \`jsonapi:"primary,${pluralize(
        ir.resourceType
      )},omitempty"\``
    );
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const primitiveType = generatePrimitiveType({
      parent: modelName,
      field,
      ir: fieldIR,
    });
    nested.push(...primitiveType.nested);
    modelStructLines.push(
      `${camelize(field)} ${
        primitiveType.output
      } \`jsonapi:"attribute" json:"${underscore(field)},omitempty"\``
    );
  }

  const struct = `
type ${modelName} struct {
  ${modelStructLines.join("\n")}
}
`;
  return [struct, ...nested].join("\n");
}

function generateClientModelFile({ ir, name }: { ir: IRType; name: string }) {
  const modelName = camelize(name);

  return `
// DO NOT MODIFY: This file is generated by internal/providergen/index.ts. Any changes will be overwritten during the next build.
package apiclient

${generateClientModel({ ir, name: modelName })}
`;
}

function generateResourceIR({ swagger, name }: { swagger: any; name: string }) {
  const resourceSchema = swagger.components.schemas[name];
  if (!resourceSchema) {
    throw new Error(`Resource ${name} not found`);
  }

  const newResourceSchema = swagger.components.schemas[`new_${name}`];
  if (!newResourceSchema) {
    throw new Error(`New resource ${name} not found`);
  }

  const collectionSchema = Object.entries(
    swagger.paths as Record<string, any>
  ).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId === `list${camelize(pluralize(name))}`
  )?.[1];
  if (!collectionSchema) {
    throw new Error(`List path for ${name} not found`);
  }

  const getSchema = Object.entries(swagger.paths as Record<string, any>).find(
    ([_, pathSchema]) =>
      pathSchema.get &&
      pathSchema.get.operationId === `get${camelize(singularize(name))}`
  )?.[1]?.get;
  if (!getSchema) {
    throw new Error(`Get path for ${name} not found`);
  }

  // Get path ID parameter
  const pathIdParameter = collectionSchema?.parameters?.[0]?.name as
    | string
    | undefined;
  const pathIdIR = pathIdParameter
    ? toIR({
        schema: resourceSchema.properties[pathIdParameter],
        required: null,
      })
    : null;

  const getHasQueryParams =
    getSchema?.parameters?.some((param) => param.in === "query") ?? false;

  // Generate immediate representation of the resource
  const irFields = toIR({
    schema: resourceSchema,
    required: newResourceSchema.required,
  });
  if (irFields.kind !== "object") {
    throw new Error("Resource root must be an object");
  }

  const ir: IRResource = {
    kind: "resource",
    resourceType: name,
    listPathIdParam:
      pathIdParameter && pathIdIR
        ? { name: pathIdParameter, element: pathIdIR }
        : null,
    getHasQueryParams,
    idElement: {
      kind: "string",
      computedOptionalRequired: "computed",
      description: `The ID of the ${humanize(name, true)}`,
    },
    fields: irFields.fields,
  };

  return ir;
}

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  console.log("🚀 Fetching Rootly Swagger...");
  let swagger = await getRootlySwagger();

  console.log("🚀 Modifying Rootly Swagger...");
  for (const mod of SWAGGER_MODS) {
    swagger = await mod(swagger);
  }

  await Bun.write(
    new URL("swagger.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  for (const name of RESOURCES) {
    const ir = generateResourceIR({ swagger, name });

    // Client
    {
      const code = generateClientModelFile({ ir, name });

      await writeAndFormatGoFile(
        new URL(`../apiclient/model_${name}_gen.go`, import.meta.url),
        code
      );
    }

    // Resource
    {
      const code = generateResource({ ir, name });

      await writeAndFormatGoFile(
        new URL(`../provider/resource_${name}_gen.go`, import.meta.url),
        code
      );
    }
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
