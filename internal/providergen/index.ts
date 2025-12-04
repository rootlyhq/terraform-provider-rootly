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

function generateGoType({
  mode,
  parent,
  field,
  ir,
}: {
  mode: "terraform" | "primitive" | "terraform_valuer";
  parent: string;
  field: string;
  ir: IRType;
}): {
  output: string;
  nested: string[];
} {
  return match([mode, ir])
    .with(["terraform", { kind: "string" }], () => ({
      output: "types.String",
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
    .with(["primitive", { kind: "string" }], () => ({
      output: "string",
      nested: [],
    }))
    .with(["terraform", { kind: "bool" }], () => ({
      output: "types.Bool",
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
    .with(["primitive", { kind: "bool" }], () => ({
      output: "bool",
      nested: [],
    }))
    .with(["terraform", { kind: "int" }], () => ({
      output: "types.Int64",
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
      ([_, ir]) => {
        const inner = generateGoType({
          mode,
          parent: camelize(`${parent}_${field}`),
          field: "Item",
          ir: ir.element,
        });
        return {
          output: `// TODO`,
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
      // TODO: validate primitive type
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
      // TODO: validate primitive type
      const inner = generateGoType({
        mode: "primitive",
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `// TODO`,
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
    .with(["terraform_valuer", { kind: "object" }], ([_, ir]) => {
      // const structName = camelize(`${parent}_${field}`);
      // const struct = generateModel({ name: structName, ir });
      // return {
      //   output: structName,
      //   nested: [struct],
      // };
      return {
        output: `// TODO`,
        nested: [],
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
      throw new Error(`Unsupported IR type for ${mode}: ${JSON.stringify(ir)}`);
    });
}

// interface Hints {
//   requiredAttributes?: string[];
// }

// function generateResourceSchemaAttribute({
//   swagger,
//   name,
//   schema,
//   hints,
// }: {
//   swagger: any;
//   name: string;
//   schema: any;
//   hints: Hints;
// }) {
//   const isRequired = hints.requiredAttributes?.includes(name) ?? false;

//   const common = `
//     ${isRequired ? "Required: true," : "Optional: true,"}
//     ${
//       schema.description
//         ? `Description: ${JSON.stringify(schema.description)},`
//         : ""
//     }
//     ${
//       schema.description
//         ? `MarkdownDescription: ${JSON.stringify(schema.description)},`
//         : ""
//     }
//   `;

//   return match(schema)
//     .with({ type: "string" }, () => {
//       return `schema.StringAttribute{
//         ${common}
//       }`;
//     })
//     .with({ type: "boolean" }, () => {
//       return `schema.BoolAttribute{
//         ${common}
//       }`;
//     })
//     .with({ type: "integer" }, () => {
//       return `schema.Int64Attribute{
//         ${common}
//       }`;
//     })
//     .with({ type: "array", items: { type: "string" } }, () => {
//       return `schema.ListAttribute{
//         ${common}
//         CustomType: supertypes.NewListTypeOf[string](ctx),
//       }`;
//     })
//     .with({ type: "array", items: { type: "object" } }, () => {
//       return `schema.ListNestedAttribute{
//         ${common}
//         NestedObject: schema.NestedAttributeObject{
//           Attributes: map[string]schema.Attribute{
//             ${Object.entries(schema.items.properties)
//               .map(([name, propertySchema]) => {
//                 const hints: Hints = {
//                   requiredAttributes: schema.items.required,
//                 };
//                 try {
//                   return `"${name}": ${generateResourceSchemaAttribute({
//                     swagger,
//                     name,
//                     schema: propertySchema,
//                     hints,
//                   })},`;
//                 } catch (e) {
//                   console.error(e);
//                   return `// TODO: Unsupported type ${name}: ${schema.type}`;
//                 }
//               })
//               .join("\n")}
//           },
//         },
//       }`;
//     })
//     .otherwise(() => {
//       throw new Error(`Unsupported type ${schema.type}`);
//     });
// }

// function generateResourceSchemav1({
//   swagger,
//   name,
// }: {
//   swagger: any;
//   name: string;
// }) {
//   const resourceSchema = swagger.components.schemas[name];
//   if (!resourceSchema) {
//     throw new Error(`Resource ${name} not found`);
//   }

//   const newResourceSchema = swagger.components.schemas[`new_${name}`];
//   if (!newResourceSchema) {
//     throw new Error(`New resource ${name} not found`);
//   }

//   const requiredAttributes =
//     newResourceSchema.properties.data.properties.attributes.required;
//   const hints: Hints = {
//     requiredAttributes,
//   };

//   return `
// schema.Schema{
//   Attributes: map[string]schema.Attribute{
//     ${Object.entries(resourceSchema.properties)
//       .map(([name, schema]) => {
//         try {
//           return `"${underscore(name)}": ${generateResourceSchemaAttribute({
//             swagger,
//             name,
//             schema,
//             hints,
//           })},`;
//         } catch (e) {
//           console.error(e);
//           return `// TODO: Unsupported type ${name}: ${schema.type}`;
//         }
//       })
//       .join("\n")}
//   },
// }
// `;
// }

// function generateModelType({
//   parentName,
//   fieldName,
//   schema,
// }: {
//   parentName: string;
//   fieldName: string;
//   schema: any;
// }): {
//   type: string;
//   nested: string[];
// } {
//   return match(schema)
//     .with({ type: "string" }, () => ({
//       type: "types.String",
//       nested: [],
//     }))
//     .with({ type: "boolean" }, () => ({
//       type: "types.Bool",
//       nested: [],
//     }))
//     .with({ type: "integer" }, () => ({
//       type: "types.Int64",
//       nested: [],
//     }))
//     .with({ type: "array", items: { type: "string" } }, () => ({
//       type: "supertypes.ListValueOf[string]",
//       nested: [],
//     }))
//     .with({ type: "array", items: { type: "object" } }, ({ items }) => {
//       const nestedName = camelize(`${parentName}_${fieldName}_item`);
//       const struct = generateModelStruct({
//         name: nestedName,
//         schema: items,
//       });

//       return {
//         type: `supertypes.ListNestedObjectValueOf[${nestedName}]`,
//         nested: [struct],
//       };
//     })
//     .with({ type: "object" }, () => {
//       const nestedName = camelize(`${parentName}_${fieldName}`);
//       const struct = generateModelStruct({
//         name: nestedName,
//         schema,
//       });

//       return {
//         type: nestedName,
//         nested: [struct],
//       };
//     })

//     .otherwise(() => {
//       throw new Error(`Unsupported type in recursive model: ${schema.type}`);
//     });
// }

// function generateModelStruct({ name, schema }: { name: string; schema: any }) {
//   const props = schema.properties ?? [];
//   const lines: string[] = [];
//   const extraStructs: string[] = [];

//   for (const [fieldName, fieldSchema] of Object.entries(props)) {
//     const { type, nested } = generateModelType({
//       parentName: name,
//       fieldName,
//       schema: fieldSchema,
//     });

//     lines.push(
//       `${camelize(fieldName)} ${type} \`tfsdk:"${underscore(fieldName)}"\``
//     );

//     extraStructs.push(...nested);
//   }

//   const struct = `
// type ${name} struct {
//   ${lines.join("\n  ")}
// }
// `;

//   return [struct, ...extraStructs].join("\n");
// }

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
  ${toClientModelLines.join("\n  ")}
  return &out, nil
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

  const common = `
    ${ir.computedOptionalRequired === "required" ? "Required: true," : ""}
    ${ir.computedOptionalRequired === "optional" ? "Optional: true," : ""}
    ${ir.computedOptionalRequired === "computed" ? "Computed: true," : ""}
    ${ir.description ? `Description: ${JSON.stringify(ir.description)},` : ""}
    ${
      ir.description
        ? `MarkdownDescription: ${JSON.stringify(ir.description)},`
        : ""
    }
  `;
  return (
    match(ir)
      .with({ kind: "string" }, () => {
        return `schema.StringAttribute{
          ${common}
        }`;
      })
      .with({ kind: "bool" }, () => {
        return `schema.BoolAttribute{
          ${common}
        }`;
      })
      .with({ kind: "int" }, () => {
        return `schema.Int64Attribute{
          ${common}
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
`;
}

function generateClientModel({ ir, name }: { ir: IRType; name: string }) {
  if (ir.kind !== "resource" && ir.kind !== "object") {
    throw new Error("Model root must be a resource");
  }

  const lines: string[] = [];
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
    lines.push(
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
    lines.push(
      `${camelize(field)} ${type} \`jsonapi:"attribute" json:"${underscore(
        field
      )},omitempty"\``
    );
  }

  const struct = `
type ${modelName} struct {
  ${lines.join("\n")}
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
