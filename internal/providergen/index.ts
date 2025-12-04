import { camelize, humanize, pluralize, underscore } from "inflection";
import { match, P } from "ts-pattern";
import { toIR, type IRResource, type IRType } from "./ir";

const RESOURCES = ["alert_route"];

async function getRootlySwagger() {
  const response = await fetch(
    "https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json"
  );
  return (await response.json()) as any;
}

interface GenerateGoTypeResult {
  output: string;
  nested: string[];
}

function generateGoType({
  mode,
  parent,
  field,
  ir,
}: {
  mode: "terraform" | "primitive" | "terraform_valuer" | "fill_model";
  parent: string;
  field: string;
  ir: IRType;
}): GenerateGoTypeResult {
  console.log({ mode, parent, field, ir });
  return match([mode, ir])
    .returnType<GenerateGoTypeResult>()
    .with(["terraform", { kind: "string" }], () => ({
      output: "supertypes.StringValue",
      nested: [],
    }))
    .with(["terraform_valuer", { kind: "string" }], () => ({
      output: `
        if !m.${camelize(field)}.IsNull() {
          out.${camelize(field)} = m.${camelize(field)}.ValueString()
        }
      `,
      nested: [],
    }))
    .with(["fill_model", { kind: "string" }], () => ({
      output: `out.${camelize(
        field
      )} = supertypes.NewStringValueOrNull(in.${camelize(field)})`,
      nested: [],
    }))
    .with(["primitive", { kind: "string" }], () => ({
      output: "string",
      nested: [],
    }))
    .with(["terraform", { kind: "bool" }], () => ({
      output: "supertypes.BoolValue",
      nested: [],
    }))
    .with(["terraform_valuer", { kind: "bool" }], () => ({
      output: `
        if !m.${camelize(field)}.IsNull() {
          out.${camelize(field)} = m.${camelize(field)}.ValueBool()
        }
      `,
      nested: [],
    }))
    .with(["fill_model", { kind: "bool" }], () => ({
      output: `out.${camelize(field)} = supertypes.NewBoolValue(in.${camelize(
        field
      )})`,
      nested: [],
    }))
    .with(["primitive", { kind: "bool" }], () => ({
      output: "bool",
      nested: [],
    }))
    .with(["terraform", { kind: "int" }], () => ({
      output: "supertypes.Int64Value",
      nested: [],
    }))
    .with(["terraform_valuer", { kind: "int" }], () => ({
      output: `
        if !m.${camelize(field)}.IsNull() {
          out.${camelize(field)} = m.${camelize(field)}.ValueInt64()
        }
      `,
      nested: [],
    }))
    .with(["fill_model", { kind: "int" }], () => ({
      output: `out.${camelize(field)} = supertypes.NewInt64Value(in.${camelize(
        field
      )})`,
      nested: [],
    }))
    .with(["primitive", { kind: "int" }], () => ({
      output: "int64",
      nested: [],
    }))
    .with(
      ["terraform", { kind: "array", element: { kind: "object" } }],
      ([_, ir]) => {
        const inner = generateGoType({
          mode,
          parent: camelize(`${parent}_${field}`),
          field: "Item",
          ir: ir.element,
        });
        return {
          output: `supertypes.ListNestedObjectValueOf[${inner.output}]`,
          nested: inner.nested,
        };
      }
    )
    .with(
      ["terraform_valuer", { kind: "array", element: { kind: "object" } }],
      () => {
        return {
          output: `
            if !m.${camelize(field)}.IsNull() {
              for _, item := range m.${camelize(field)}.MustGet(ctx) {
                itemClientModel, err := item.ToClientModel(ctx)
                if err != nil {
                  return nil, err
                }
                out.${camelize(field)} = append(out.${camelize(
            field
          )}, *itemClientModel)			
              }
            }
          `,
          nested: [],
        };
      }
    )
    .with(
      ["fill_model", { kind: "array", element: { kind: "object" } }],
      ([_, ir]) => {
        const inner = generateGoType({
          mode,
          parent: camelize(`${parent}_${field}`),
          field: "Item",
          ir: ir.element,
        });
        return {
          output: `
            {
              var elements []${camelize(`${parent}_${field}`)}Item
              for _, item := range in.${camelize(field)} {
                var element ${camelize(`${parent}_${field}`)}Item
                err := Fill${camelize(
                  `${parent}_${field}`
                )}Item(ctx, item, &element)
                if err != nil {
                  return err
                }
                elements = append(elements, element)
              }
              out.${camelize(
                field
              )} = supertypes.NewListNestedObjectValueOfValueSlice(ctx, elements)
            }
          `,
          nested: inner.nested,
        };
      }
    )
    .with(
      ["primitive", { kind: "array", element: { kind: "object" } }],
      ([_, ir]) => {
        const inner = generateGoType({
          mode,
          parent: camelize(`${parent}_${field}`),
          field: "Item",
          ir: ir.element,
        });
        return {
          output: `[]${inner.output}`,
          nested: inner.nested,
        };
      }
    )
    .with(["terraform", { kind: "array" }], ([_, ir]) => {
      const inner = generateGoType({
        mode: "primitive",
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `supertypes.ListValueOf[${inner.output}]`,
        nested: inner.nested,
      };
    })
    .with(["terraform_valuer", { kind: "array" }], ([_, ir]) => {
      const inner = generateGoType({
        mode: "primitive",
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `
          if !m.${camelize(field)}.IsNull() {
            out.${camelize(field)} = m.${camelize(field)}.MustGet(ctx)
          }
        `,
        nested: inner.nested,
      };
    })
    .with(["fill_model", { kind: "array" }], ([_, ir]) => {
      const inner = generateGoType({
        mode,
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `out.${camelize(
          field
        )} = supertypes.NewListValueOfSlice(ctx, in.${camelize(field)})`,
        nested: inner.nested,
      };
    })
    .with(["primitive", { kind: "array" }], ([_, ir]) => {
      const inner = generateGoType({
        mode,
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `[]${inner.output}`,
        nested: inner.nested,
      };
    })
    .with(["terraform", { kind: "object" }], ([_, ir]) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateModel({ name: structName, ir });
      return {
        output: structName,
        nested: [struct],
      };
    })
    .with(["fill_model", { kind: "object" }], ([_, ir]) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateFillModel({ name: structName, ir });
      return {
        output: structName,
        nested: [struct],
      };
    })
    .with(["primitive", { kind: "object" }], ([_, ir]) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateClientModel({ name: structName, ir });
      return {
        output: structName,
        nested: [struct],
      };
    })
    .otherwise(() => {
      throw new Error(
        `Unsupported IR type for field ${parent}.${field}, mode ${mode}: ${JSON.stringify(
          ir
        )}`
      );
    });
}

function generateModel({ name, ir }: { name: string; ir: IRType }) {
  if (ir.kind !== "object" && ir.kind !== "resource") {
    throw new Error("Model root must be an object");
  }

  const modelStructLines: string[] = [];
  const toClientModelLines: string[] = [];
  const nested: string[] = [];

  if (ir.kind === "resource") {
    const { output } = generateGoType({
      mode: "terraform",
      parent: name,
      field: "id",
      ir: ir.idElement,
    });
    modelStructLines.push(`Id ${output} \`tfsdk:"id"\``);
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const terraformType = generateGoType({
      mode: "terraform",
      parent: name,
      field,
      ir: fieldIR,
    });
    const terraformValuerType = generateGoType({
      mode: "terraform_valuer",
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
    toClientModelLines.push(terraformValuerType.output);
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
    const { output } = generateGoType({
      mode: "fill_model",
      parent: name,
      field: "id",
      ir: ir.idElement,
    });
    fillModelLines.push(output);
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const result = generateGoType({
      mode: "fill_model",
      parent: name,
      field,
      ir: fieldIR,
    });

    nested.push(...result.nested);
    fillModelLines.push(result.output);
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
  ir: IRType;
}) {
  if (ir.kind === "resource") {
    throw new Error("Resource field cannot be an attribute");
  }

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

  return (
    match(ir)
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
        const { output: type } = generateGoType({
          mode: "primitive",
          parent,
          field,
          ir: element,
        });

        return `schema.ListAttribute{
          ${common}
          CustomType: supertypes.NewListTypeOf[${type}](ctx),
        }`;
      })
      // .with({ kind: "object" }, () => {
      //   const nestedName = camelize(`${parentName}_${fieldName}`);
      //   const struct = generateModelStruct({
      //     name: nestedName,
      //     schema,
      //   });

      //   return `${nestedName}`;
      // })
      .otherwise(() => {
        throw new Error(`Unsupported IR type: ${JSON.stringify(ir)}`);
      })
  );
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

function generateResource({ ir, name }: { ir: IRType; name: string }) {
  const resourceName = `${camelize(name)}Resource`;
  const modelName = `${camelize(name)}Model`;

  return `
// DO NOT MODIFY: This file is generated by internal/providergen/index.ts. Any changes will be overwritten during the next build.
package provider

import (
  "context"
  "github.com/hashicorp/terraform-plugin-framework/resource"
  "github.com/hashicorp/terraform-plugin-framework/resource/schema"
  "github.com/hashicorp/terraform-plugin-framework/types"
  "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/apiclient"
)

var _ resource.Resource = (*${resourceName})(nil)
var _ resource.ResourceWithImportState = (*${resourceName})(nil)

func New${resourceName}() resource.Resource {
  return &${resourceName}{}
}

type ${resourceName} struct {
  baseResource
  extendableResource[${modelName}]
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
  r.create(ctx, &data, resp)

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
  r.read(ctx, &data, resp)

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
  r.update(ctx, plan, &data, resp)

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
  r.delete(ctx, &data, resp)
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
    const { output: type } = generateGoType({
      mode: "primitive",
      parent: modelName,
      field: "ir",
      ir: ir.idElement,
    });
    modelStructLines.push(
      `Id ${type} \`jsonapi:"primary,${pluralize(ir.resourceType)},omitempty"\``
    );
  }

  for (const [field, fieldIR] of Object.entries(ir.fields)) {
    const { output: type, nested: nx } = generateGoType({
      mode: "primitive",
      parent: modelName,
      field,
      ir: fieldIR,
    });
    nested.push(...nx);
    modelStructLines.push(
      `${camelize(field)} ${type} \`jsonapi:"attribute" json:"${underscore(
        field
      )},omitempty"\``
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

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  console.log("🚀 Fetching Rootly Swagger...");
  const swagger = await getRootlySwagger();

  await Bun.write(
    new URL("swagger.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  for (const name of RESOURCES) {
    const resourceSchema = swagger.components.schemas[name];
    if (!resourceSchema) {
      throw new Error(`Resource ${name} not found`);
    }

    const newResourceSchema = swagger.components.schemas[`new_${name}`];
    if (!newResourceSchema) {
      throw new Error(`New resource ${name} not found`);
    }

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
      idElement: {
        kind: "string",
        computedOptionalRequired: "computed",
        description: `The ID of the ${humanize(name, true)}`,
      },
      fields: irFields.fields,
    };

    console.log(JSON.stringify(ir, null, 2));

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
