import { camelize, pluralize, underscore } from "inflection";
import { match } from "ts-pattern";
import {
  generateResourceIR,
  type IRObject,
  type IRResource,
  type IRType,
} from "./ir";
import { SWAGGER_MODS } from "./swagger-mods";
import { IR_MODS } from "./ir-mods";
import { RESOURCES } from "./settings";

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
  hasDiags: boolean;
} {
  return match(ir)
    .returnType<{ output: string; hasDiags: boolean }>()
    .with({ kind: "string", nullable: false }, () => ({
      output: `${prefix}${camelize(field)}.ValueString()`,
      hasDiags: false,
    }))
    .with({ kind: "string", nullable: true }, () => ({
      output: `${prefix}${camelize(field)}.ValueStringPointer()`,
      hasDiags: false,
    }))
    .with({ kind: "bool", nullable: false }, () => ({
      output: `${prefix}${camelize(field)}.ValueBool()`,
      hasDiags: false,
    }))
    .with({ kind: "bool", nullable: true }, () => ({
      output: `${prefix}${camelize(field)}.ValueBoolPointer()`,
      hasDiags: false,
    }))
    .with({ kind: "int", nullable: false }, () => ({
      output: `${prefix}${camelize(field)}.ValueInt64()`,
      hasDiags: false,
    }))
    .with({ kind: "int", nullable: true }, () => ({
      output: `${prefix}${camelize(field)}.ValueInt64Pointer()`,
      hasDiags: false,
    }))
    .with(
      { kind: "array", element: { kind: "object" }, nullable: false },
      () => {
        const structName = camelize(`${parent}_${field}_item`);
        const fieldName = `${prefix}${camelize(field)}`;
        return {
          output: `
          func() ([]apiclient.${structName}, diag.Diagnostics) {
            var diags diag.Diagnostics
            var itemClientModels []apiclient.${structName}
            for _, item := range ${fieldName}.DiagsGet(ctx, diags) {
              itemClientModel := tfutils.MergeDiagnostics(item.ToClientModel(ctx))(&diags)
              itemClientModels = append(itemClientModels, *itemClientModel)
            }
            if diags.HasError() {
              return nil, diags
            }
            return itemClientModels, nil
          }()
        `.trim(),
          hasDiags: true,
        };
      }
    )
    .with(
      { kind: "array", element: { kind: "object" }, nullable: true },
      () => {
        const structName = camelize(`${parent}_${field}_item`);
        const fieldName = `${prefix}${camelize(field)}`;
        return {
          output: `
          func() (*[]apiclient.${structName}, diag.Diagnostics) {
            var diags diag.Diagnostics
            var itemClientModels []apiclient.${structName}
            for _, item := range ${fieldName}.DiagsGet(ctx, diags) {
              itemClientModel := tfutils.MergeDiagnostics(item.ToClientModel(ctx))(&diags)
              itemClientModels = append(itemClientModels, *itemClientModel)
            }
            if diags.HasError() {
              return nil, diags
            }
            return &itemClientModels, nil
          }()
        `.trim(),
          hasDiags: true,
        };
      }
    )
    .with({ kind: "array", nullable: false }, () => ({
      output: `${prefix}${camelize(field)}.Get(ctx)`,
      hasDiags: true,
    }))
    .with({ kind: "array", nullable: true }, (ir) => {
      const inner = generatePrimitiveType({
        parent,
        field,
        ir,
      });
      return {
        output: `
          func() (${inner.output}, diag.Diagnostics) {
            v, diags := ${prefix}${camelize(field)}.Get(ctx)
            if diags.HasError() {
              return nil, diags
            }
            return ptr.Ptr(v), nil
          }()
        `.trim(),
        hasDiags: true,
      };
    })
    .with({ kind: "object", nullable: false }, () => ({
      output: `
        func() (apiclient.${camelize(`${parent}_${field}`)}, diag.Diagnostics) {
          var diags diag.Diagnostics
          model := ${prefix}${camelize(field)}.DiagsGet(ctx, diags)
          if diags.HasError() {
            return apiclient.${camelize(`${parent}_${field}`)}{}, diags
          }
          clientModel := tfutils.MergeDiagnostics(model.ToClientModel(ctx))(&diags)
          if diags.HasError() {
            return apiclient.${camelize(`${parent}_${field}`)}{}, diags
          }
          return *clientModel, nil
        }()
      `.trim(),
      hasDiags: true,
    }))
    .with({ kind: "object", nullable: true }, () => ({
      output: `
        func() (*apiclient.${camelize(
          `${parent}_${field}`
        )}, diag.Diagnostics) {
          var diags diag.Diagnostics
          model := ${prefix}${camelize(field)}.DiagsGet(ctx, diags)
          if diags.HasError() {
            return nil, diags
          }
          clientModel := tfutils.MergeDiagnostics(model.ToClientModel(ctx))(&diags)
          if diags.HasError() {
            return nil, diags
          }
          return clientModel, nil
        }()
      `.trim(),
      hasDiags: true,
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
      const fieldType = ir.distinct ? "Set" : "List";
      const inner = generateTerraformType({
        inObject: true,
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      return {
        output: `supertypes.${fieldType}NestedObjectValueOf[${inner.output}]`,
        nested: inner.nested,
      };
    })
    .with({ kind: "array" }, (ir) => {
      const fieldType = ir.distinct ? "Set" : "List";
      const inner = generatePrimitiveType({
        parent,
        field,
        ir: ir.element,
      });
      return {
        output: `supertypes.${fieldType}ValueOf[${inner.output}]`,
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
  const result = match(ir)
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
        parent,
        field,
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

  return {
    output: ir.nullable ? `*${result.output}` : result.output,
    nested: result.nested,
  };
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
  hasDiags: boolean;
} {
  return match(ir)
    .returnType<{
      output: string;
      nested: string[];
      hasDiags: boolean;
    }>()
    .with({ kind: "string", nullable: false }, () => ({
      output: `supertypes.NewStringValueOrNull(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "string", nullable: true }, () => ({
      output: `supertypes.NewStringPointerValueOrNull(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "bool", nullable: false }, () => ({
      output: `supertypes.NewBoolValue(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "bool", nullable: true }, () => ({
      output: `supertypes.NewBoolPointerValueOrNull(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "int", nullable: false }, () => ({
      output: `supertypes.NewInt64Value(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "int", nullable: true }, () => ({
      output: `supertypes.NewInt64PointerValueOrNull(in.${camelize(field)})`,
      nested: [],
      hasDiags: false,
    }))
    .with({ kind: "array", element: { kind: "object" } }, (ir) => {
      const fieldType = ir.distinct ? "Set" : "List";
      const inner = generateModelValuer({
        parent: camelize(`${parent}_${field}`),
        field: "Item",
        ir: ir.element,
      });
      const itemType = camelize(`${parent}_${field}_item`);
      return {
        output: `
          func() (supertypes.${fieldType}NestedObjectValueOf[${itemType}], diag.Diagnostics) {
            var diags diag.Diagnostics
            var elements []${itemType}
            for _, item := range in.${camelize(field)} {
              var element ${itemType}
              diags.Append(Fill${itemType}(ctx, item, &element)...)
              elements = append(elements, element)
            }

            if diags.HasError() {
              return supertypes.New${fieldType}NestedObjectValueOfNull[${itemType}](ctx), diags
            }
            return supertypes.New${fieldType}NestedObjectValueOfValueSlice(ctx, elements), nil
          }()
        `.trim(),
        nested: inner.nested,
        hasDiags: true,
      };
    })
    .with({ kind: "array", nullable: false }, (ir) => {
      const fieldType = ir.distinct ? "Set" : "List";
      const inner = generateModelValuer({
        parent,
        field,
        ir: ir.element,
      });
      return {
        output: `supertypes.New${fieldType}ValueOfSlice(ctx, in.${camelize(
          field
        )})`,
        nested: inner.nested,
        hasDiags: false,
      };
    })
    .with({ kind: "array", nullable: true }, (ir) => {
      const fieldType = ir.distinct ? "Set" : "List";
      const primitive = generatePrimitiveType({
        parent,
        field,
        ir: ir.element,
      });
      const inner = generateModelValuer({
        parent,
        field,
        ir: ir.element,
      });
      return {
        output: `
          func() supertypes.${fieldType}ValueOf[${primitive.output}] {
            if in.${camelize(field)} == nil {
              return supertypes.New${fieldType}ValueOfNull[${
          primitive.output
        }](ctx)
            }
            return supertypes.New${fieldType}ValueOfSlice(ctx, *in.${camelize(
          field
        )})
          }()
        `.trim(),
        nested: inner.nested,
        hasDiags: false,
      };
    })
    .with({ kind: "object", nullable: false }, (ir) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateFillModel({ name: structName, ir });
      return {
        output: `
          func() (supertypes.SingleNestedObjectValueOf[${structName}], diag.Diagnostics) {
            var diags diag.Diagnostics
            var out ${structName}
            diags.Append(Fill${structName}(ctx, in.${camelize(field)}, &out)...)

            if diags.HasError() {
              return supertypes.NewSingleNestedObjectValueOfNull[${structName}](ctx), diags
            }
            return supertypes.NewSingleNestedObjectValueOf(ctx, &out), nil
          }()
        `.trim(),
        nested: [struct],
        hasDiags: true,
      };
    })
    .with({ kind: "object", nullable: true }, (ir) => {
      const structName = camelize(`${parent}_${field}`);
      const struct = generateFillModel({ name: structName, ir });
      return {
        output: `
          func() (supertypes.SingleNestedObjectValueOf[${structName}], diag.Diagnostics) {
            var diags diag.Diagnostics

            if in.${camelize(field)} == nil {
              return supertypes.NewSingleNestedObjectValueOfNull[${structName}](ctx), diags
            }

            var out ${structName}
            diags.Append(Fill${structName}(ctx, *in.${camelize(
          field
        )}, &out)...)

            if diags.HasError() {
              return supertypes.NewSingleNestedObjectValueOfNull[${structName}](ctx), diags
            }
            return supertypes.NewSingleNestedObjectValueOf(ctx, &out), nil
          }()
        `.trim(),
        nested: [struct],
        hasDiags: true,
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

    toClientModelLines.push(
      `if !m.${camelize(field)}.IsNull() && !m.${camelize(field)}.IsUnknown() {`
    );
    if (terraformValuer.hasDiags) {
      toClientModelLines.push(
        `out.${camelize(field)} = tfutils.MergeDiagnostics(${
          terraformValuer.output
        })(&diags)`
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

func (m *${name}) ToClientModel(ctx context.Context) (*apiclient.${name}, diag.Diagnostics) {
  var diags diag.Diagnostics
  var out apiclient.${name}
  ${toClientModelLines.join("\n")}
  if diags.HasError() {
    return nil, diags
  }
  return &out, nil
}
`;

  return [struct, ...nested].join("\n");
}

function generateFillModel({
  name,
  ir,
}: {
  name: string;
  ir: IRResource | IRObject;
}) {
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

    if (modelValuer.hasDiags) {
      fillModelLines.push(
        `out.${camelize(field)} = tfutils.MergeDiagnostics(${
          modelValuer.output
        })(&diags)`
      );
    } else {
      fillModelLines.push(`out.${camelize(field)} = ${modelValuer.output}`);
    }
  }

  const struct = `
func Fill${name}(ctx context.Context, in apiclient.${name}, out *${name}) (diags diag.Diagnostics) {
  ${fillModelLines.join("\n")}
  return
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
  if (
    ir.computedOptionalRequired === "optional" ||
    ir.computedOptionalRequired === "computed_optional"
  ) {
    commonLines.push("Optional: true,");
  }
  if (
    ir.computedOptionalRequired === "computed" ||
    ir.computedOptionalRequired === "computed_optional"
  ) {
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
    .with({ kind: "string" }, (ir) => {
      const validators = ir.choices
        ? `Validators: []validator.String{
            stringvalidator.OneOf(${ir.choices
              .map((choice) => JSON.stringify(choice))
              .join(", ")}),
          },`
        : "";
      return `schema.StringAttribute{
          ${common}
          CustomType: supertypes.StringType{},
          ${validators}
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
      const fieldType = ir.distinct ? "Set" : "List";
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

      return `schema.${fieldType}NestedAttribute{
          ${common}
          CustomType: supertypes.New${fieldType}NestedObjectTypeOf[${structType}](ctx),
          NestedObject: schema.NestedAttributeObject{
            Attributes: map[string]schema.Attribute{
              ${attrs}
            },
          },
        }`;
    })
    .with({ kind: "array" }, (ir) => {
      const fieldType = ir.distinct ? "Set" : "List";
      const primitiveType = generatePrimitiveType({
        parent,
        field,
        ir: ir.element,
      });

      return `schema.${fieldType}Attribute{
          ${common}
          CustomType: supertypes.New${fieldType}TypeOf[${primitiveType.output}](ctx),
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

function generateResourceSchema({
  name,
  ir,
}: {
  name: string;
  ir: IRResource;
}) {
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

    if (terraformValuer.hasDiags) {
      lines.push(
        `modelIn.${camelize(field)} = tfutils.MergeDiagnostics(${
          terraformValuer.output
        })(&resp.Diagnostics)`
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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/jianyuan/go-utils/ptr"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/tfutils"
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
  modelIn := tfutils.MergeDiagnostics(data.ToClientModel(ctx))(&resp.Diagnostics)

  if resp.Diagnostics.HasError() {
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

  if resp.Diagnostics.HasError() {
    return
  }

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

function generateClientModel({
  ir,
  name,
}: {
  ir: IRResource | IRObject;
  name: string;
}) {
  const modelStructLines: string[] = [];
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

function generateClientModelFile({
  ir,
  name,
}: {
  ir: IRResource;
  name: string;
}) {
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
  let swagger = await getRootlySwagger();
  await Bun.write(
    new URL("out/swagger.original.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  console.log("🚀 Modifying Rootly Swagger...");
  for (const mod of SWAGGER_MODS) {
    swagger = await mod(swagger);
  }
  await Bun.write(
    new URL("out/swagger.modified.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  for (const name of RESOURCES) {
    let ir = generateResourceIR({ swagger, name });
    await Bun.write(
      new URL(`out/ir_${name}.original.json`, import.meta.url),
      JSON.stringify(ir, null, 2)
    );

    if (IR_MODS[name]) {
      ir = await IR_MODS[name](ir);
    }
    await Bun.write(
      new URL(`out/ir_${name}.modified.json`, import.meta.url),
      JSON.stringify(ir, null, 2)
    );

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
