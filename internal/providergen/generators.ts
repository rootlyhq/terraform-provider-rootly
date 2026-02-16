import { match, P } from "ts-pattern";
import type { IRType, ResourceIR } from "./ir";
import { camelize, pluralize, underscore } from "inflection";
import dedent from "dedent";

function generateTerraformValueType({
  mode,
  parent,
  attribute,
}: {
  mode: "legacy" | "modern";
  parent: string;
  attribute: IRType;
}) {
  return match([mode, attribute])
    .with([P.any, { type: "string" }], () => "supertypes.StringValue")
    .with([P.any, { type: "int" }], () => "supertypes.Int64Value")
    .with([P.any, { type: "bool" }], () => "supertypes.BoolValue")
    .with(
      [P.any, { type: "list", elementType: "string" }],
      () => "supertypes.ListValueOf[string]"
    )
    .with(
      [P.any, { type: "set", elementType: "string" }],
      () => "supertypes.SetValueOf[string]"
    )
    .with(
      [P.any, { type: "set_nested" }],
      () =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(attribute.name)}Item]`
    )
    .with(
      ["modern", { type: "object" }],
      () =>
        `supertypes.SingleNestedObjectValueOf[${parent}${camelize(attribute.name)}]`
    )
    .with(
      ["legacy", { type: "object" }],
      () =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(attribute.name)}]`
    )
    .exhaustive();
}

function generateResourceUpdateDiff(resource: ResourceIR) {
  const lines: string[] = [];

  for (const attribute of resource.attributes) {
    lines.push(
      `if !data.${camelize(attribute.name)}.Equal(plan.${camelize(attribute.name)}) {`
    );
    lines.push(
      generateTerraformToPrimitive({
        mode: resource.mode,
        parent: resource.name,
        attribute,
        srcVar: "plan",
        destVar: "modelIn",
        diagsVar: "resp.Diagnostics",
      })
    );
    lines.push(`}`);
  }

  return lines.join("\n");
}

function generateTerraformSchemaAttribute({
  mode,
  parent,
  attribute,
}: {
  mode: "legacy" | "modern";
  parent: string;
  attribute: IRType;
}) {
  // Description
  const descriptionParts = [];
  if (attribute.description) {
    descriptionParts.push(attribute.description);
  }
  if (attribute.type === "string" && attribute.enum) {
    descriptionParts.push(
      `Value must be one of ${attribute.enum.map((value) => `\`${value}\``).join(", ")}.`
    );
  }
  if (attribute.deprecationMessage) {
    descriptionParts.push(`**Deprecated** ${attribute.deprecationMessage}`);
  }
  const description = descriptionParts.join(" ");

  // Common parts
  const commonParts: string[] = [];
  commonParts.push(`MarkdownDescription: ${JSON.stringify(description)},`);
  if (attribute.deprecationMessage) {
    commonParts.push(
      `DeprecationMessage: ${JSON.stringify(attribute.deprecationMessage)},`
    );
  }

  // Attribute parts
  const attributeParts: string[] = [];
  attributeParts.push(
    match(attribute.computedOptionalRequired)
      .with("required", () => "Required: true,")
      .with("computed", () => "Computed: true,")
      .with("computed_optional", () => "Optional: true,\nComputed: true,")
      .with("optional", () => "Optional: true,")
      .exhaustive()
  );
  if (attribute.sensitive) {
    attributeParts.push("Sensitive: true,");
  }

  const validators: string[] = [];
  if (attribute.validators) {
    validators.push(...attribute.validators);
  }
  if (attribute.type === "string" && attribute.enum) {
    validators.push(
      `stringvalidator.OneOf(${attribute.enum.map((value) => JSON.stringify(value)).join(", ")})`
    );
  }

  return match(attribute)
    .with({ type: "string" }, () => {
      const parts: string[] = [];
      parts.push("schema.StringAttribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push("CustomType: supertypes.StringType{},");
      if (validators.length > 0) {
        parts.push("Validators: []validator.String{");
        parts.push(...validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.String{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "int" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.Int64Attribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push("CustomType: supertypes.Int64Type{},");
      if (validators.length > 0) {
        parts.push("Validators: []validator.Int64{");
        parts.push(...validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Int64{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "bool" }, () => {
      const parts: string[] = [];
      parts.push("schema.BoolAttribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push("CustomType: supertypes.BoolType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Bool{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Bool{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "list", elementType: "string" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.ListAttribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push("CustomType: supertypes.NewListTypeOf[string](ctx),");
      if (validators.length > 0) {
        parts.push("Validators: []validator.List{");
        parts.push(...validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.List{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "set", elementType: "string" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.SetAttribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push("CustomType: supertypes.NewSetTypeOf[string](ctx),");
      if (validators.length > 0) {
        parts.push("Validators: []validator.Set{");
        parts.push(...validators.map((validator) => `${validator},`));
        parts.push("},");
      }
      if (attribute.planModifiers) {
        parts.push("PlanModifiers: []planmodifier.Set{");
        parts.push(
          ...attribute.planModifiers.map((modifier) => `${modifier},`)
        );
        parts.push("},");
      }
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "set_nested" }, (attribute) => {
      const parts: string[] = [];
      parts.push("schema.SetNestedAttribute{");
      parts.push(...commonParts);
      parts.push(...attributeParts);
      parts.push(
        `CustomType: supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(
          attribute.name
        )}Item](ctx),`
      );
      parts.push("NestedObject: schema.NestedAttributeObject{");
      parts.push(
        generateTerraformSchemaAttributesBlocks({
          mode,
          name: `${parent}${camelize(attribute.name)}Item`,
          attributes: attribute.attributes,
        })
      );
      parts.push("},");
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "object" }, (attribute) => {
      if (mode === "modern") {
        const parts: string[] = [];

        parts.push("schema.SingleNestedAttribute{");
        parts.push(...commonParts);
        parts.push(...attributeParts);
        parts.push(
          `CustomType: supertypes.NewSingleNestedObjectTypeOf[${parent}${camelize(attribute.name)}](ctx),`
        );
        parts.push(
          generateTerraformSchemaAttributesBlocks({
            mode,
            name: `${parent}${camelize(attribute.name)}Item`,
            attributes: attribute.attributes,
          })
        );
        parts.push("}");

        return parts.join("\n");
      }

      const parts: string[] = [];

      parts.push("schema.SetNestedBlock{");
      parts.push(...commonParts);
      parts.push(
        `CustomType: supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(attribute.name)}](ctx),`
      );
      parts.push("NestedObject: schema.NestedBlockObject{");
      parts.push(
        generateTerraformSchemaAttributesBlocks({
          mode,
          name: `${parent}${camelize(attribute.name)}Item`,
          attributes: attribute.attributes,
        })
      );
      parts.push("},");
      parts.push("Validators: []validator.Set{");
      parts.push("setvalidator.SizeAtMost(1),");
      parts.push("},");
      parts.push("}");

      return parts.join("\n");
    })
    .exhaustive();
}

function generateTerraformSchemaAttributesBlocks({
  mode,
  name,
  attributes,
}: {
  mode: "legacy" | "modern";
  name: string;
  attributes: IRType[];
}) {
  const attributeLines: string[] = [];
  const blockLines: string[] = [];

  for (const attribute of attributes) {
    match([mode, attribute])
      .with(["legacy", { type: "object" }], () => {
        blockLines.push(
          `"${attribute.name}": ${generateTerraformSchemaAttribute({
            mode,
            parent: name,
            attribute,
          })},`
        );
      })
      .otherwise(() => {
        attributeLines.push(
          `"${attribute.name}": ${generateTerraformSchemaAttribute({
            mode,
            parent: name,
            attribute,
          })},`
        );
      });
  }

  const parts: string[] = [];

  if (attributeLines.length > 0) {
    parts.push("Attributes: map[string]schema.Attribute{");
    parts.push(attributeLines.join("\n"));
    parts.push("},");
  }

  if (blockLines.length > 0) {
    parts.push("Blocks: map[string]schema.Block{");
    parts.push(blockLines.join("\n"));
    parts.push("},");
  }

  return parts.join("\n");
}

function generatePrimitiveToTerraform({
  mode,
  parent,
  attribute,
  srcVar,
  destVar,
}: {
  mode: "legacy" | "modern";
  parent: string;
  attribute: IRType;
  srcVar: string;
  destVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  const destVarName = `${destVar}.${camelize(attribute.name)}`;
  return match([mode, attribute])
    .with(
      [P.any, { type: "string", nullable: true }],
      () => `${destVarName} = supertypes.NewStringPointerValue(${srcVarName})`
    )
    .with(
      [P.any, { type: "string", sourceType: "time" }],
      () => `${destVarName} = supertypes.NewStringValue(${srcVarName}.String())`
    )
    .with(
      [P.any, { type: "string" }],
      () => `${destVarName} = supertypes.NewStringValue(${srcVarName})`
    )
    .with(
      [P.any, { type: "int" }],
      () => `${destVarName} = supertypes.NewInt64Value(${srcVarName})`
    )
    .with(
      [P.any, { type: "bool" }],
      () => `${destVarName} = supertypes.NewBoolValue(${srcVarName})`
    )
    .with(
      [P.any, { type: "list", elementType: "string" }],
      () =>
        `${destVarName} = supertypes.NewListValueOfSlice(ctx, ${srcVarName})`
    )
    .with(
      [P.any, { type: "set", nullable: true, elementType: "string" }],
      () => dedent`
        if ${srcVarName} == nil {
          ${destVarName} = supertypes.NewSetValueOfNull[string](ctx)
        } else {
          ${destVarName} = supertypes.NewSetValueOfSlice(ctx, *${srcVarName})
        }
      `
    )
    .with(
      [P.any, { type: "set", elementType: "string" }],
      () => `${destVarName} = supertypes.NewSetValueOfSlice(ctx, ${srcVarName})`
    )
    .with(
      [P.any, { type: "set_nested" }],
      (attribute) =>
        `${destVarName} = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, sliceutils.Map(func(item apiclient.${attribute.model}) ${parent}${camelize(attribute.name)}Item {
          var model ${parent}${camelize(attribute.name)}Item
          diags.Append(model.FromApi(ctx, item)...)
          return model
        }, ${srcVarName}))`
    )
    .with(
      ["legacy", { type: "object" }],
      () => dedent`
        if ${destVarName}.IsKnown() {
          if ${srcVarName} == nil {
            ${destVarName} = supertypes.NewSetNestedObjectValueOfNull[${parent}${camelize(attribute.name)}](ctx)
          } else {
            var model ${parent}${camelize(attribute.name)}
            diags.Append(model.FromApi(ctx, *${srcVarName})...)
            ${destVarName} = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, []${parent}${camelize(attribute.name)}{model})
          }
        }
      `
    )
    .with(
      ["modern", { type: "object" }],
      () => dedent`
        if ${srcVarName} == nil {
          ${destVarName} = supertypes.NewSingleNestedObjectValueOfNull[${parent}${camelize(attribute.name)}](ctx)
        } else {
          var model ${parent}${camelize(attribute.name)}
          diags.Append(model.FromApi(ctx, *${srcVarName})...)
          ${destVarName} = supertypes.NewSingleNestedObjectValueOf[${parent}${camelize(attribute.name)}](ctx, &model)
        }
      `
    )
    .exhaustive();
}

function generateTerraformToPrimitive({
  mode,
  parent,
  attribute,
  srcVar,
  destVar,
  diagsVar,
}: {
  mode: "legacy" | "modern";
  parent: string;
  attribute: IRType;
  srcVar: string;
  destVar: string;
  diagsVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  const destVarName = `${destVar}.${camelize(attribute.name)}`;
  return match([mode, attribute])
    .with(
      [P.any, { type: "string", nullable: true }],
      () => `${destVarName} = ${srcVarName}.ValueStringPointer()`
    )
    .with(
      [P.any, { type: "string" }],
      () => `${destVarName} = ${srcVarName}.ValueString()`
    )
    .with(
      [P.any, { type: "int", nullable: true }],
      () => `${destVarName} = ${srcVarName}.ValueInt64Pointer()`
    )
    .with(
      [P.any, { type: "int" }],
      () => `${destVarName} = ${srcVarName}.ValueInt64()`
    )
    .with(
      [P.any, { type: "bool", nullable: true }],
      () => `${destVarName} = ${srcVarName}.ValueBoolPointer()`
    )
    .with(
      [P.any, { type: "bool" }],
      () => `${destVarName} = ${srcVarName}.ValueBool()`
    )
    .with(
      [P.any, { type: "list", nullable: true, elementType: "string" }],
      () =>
        dedent`
          if ${srcVarName}.IsKnown() {
            ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar}))
          } else {
            ${destVarName} = nil
          }
        `
    )
    .with(
      [P.any, { type: "list", elementType: "string" }],
      () =>
        dedent`
          if ${srcVarName}.IsKnown() {
            ${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})
          } else {
            ${destVarName} = nil
          }
        `
    )
    .with(
      [P.any, { type: "set", nullable: true, elementType: "string" }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar}))
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      [P.any, { type: "set", elementType: "string" }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          ${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      [P.any, { type: "set_nested" }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar}))
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      ["modern", { type: "object", nullable: true }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          model := tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})
          if ${diagsVar}.HasError() {
            return
          }
          ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(model.ToApi(ctx))(&${diagsVar}))
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      ["modern", { type: "object" }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          model := tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})
          if ${diagsVar}.HasError() {
            return
          }
          ${destVarName} = tfutils.MergeDiagnostics(model.ToApi(ctx))(&${diagsVar})
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      ["legacy", { type: "object", nullable: true }],
      () => dedent`
        if ${srcVarName}.IsKnown() {
          models := tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})
          if ${diagsVar}.HasError() {
            return
          }

          if len(models) == 1 {
            ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(models[0].ToApi(ctx))(&${diagsVar}))
          } else {
            ${destVarName} = nil
          }
        } else {
          ${destVarName} = nil
        }
      `
    )
    .otherwise((schema) => `// TODO ${JSON.stringify(schema)}`);
}

function generateModel({
  mode,
  name,
  attributes,
}: {
  mode: "legacy" | "modern";
  name: string;
  attributes: IRType[];
}) {
  const lines: string[] = [];
  const fromApiLines: string[] = [];
  const toApiLines: string[] = [];
  const extras: string[] = [];

  for (const attribute of attributes) {
    lines.push(
      `${camelize(attribute.name)} ${generateTerraformValueType({
        mode,
        parent: name,
        attribute,
      })} \`tfsdk:"${attribute.name}"\``
    );

    fromApiLines.push(
      `${generatePrimitiveToTerraform({
        mode,
        parent: name,
        attribute,
        srcVar: "data",
        destVar: "m",
      })}${attribute.deprecationMessage ? " // Deprecated" : ""}`
    );

    toApiLines.push(
      `${generateTerraformToPrimitive({
        mode,
        parent: name,
        attribute,
        srcVar: "m",
        destVar: "data",
        diagsVar: "diags",
      })}${attribute.deprecationMessage ? " // Deprecated" : ""}`
    );

    extras.push(
      ...match(attribute)
        .with({ type: "object" }, (attribute) => [
          generateModel({
            mode,
            name: `${name}${camelize(attribute.name)}`,
            attributes: attribute.attributes,
          }),
        ])
        .otherwise(() => [])
    );
  }

  return dedent`
type ${name} struct {
  ${lines.join("\n")}
}

func (m *${name}) FromApi(ctx context.Context, data apiclient.${name}) (diags diag.Diagnostics) {
  ${fromApiLines.join("\n")}
  return
}

func (m *${name}) ToApi(ctx context.Context) (data apiclient.${name}, diags diag.Diagnostics) {
  ${toApiLines.join("\n")}
  return
}

${extras.join("\n\n")}
`;
}

function generateResourceModel(resource: ResourceIR) {
  return generateModel({
    mode: resource.mode,
    name: `${camelize(resource.name)}ResourceModel`,
    attributes: [
      ...(resource.idAttribute ? [resource.idAttribute] : []),
      ...resource.attributes,
    ],
  });
}

export function generateResource(resource: ResourceIR) {
  const baseName = camelize(resource.name);
  const resourceName = `${baseName}Resource`;
  const modelName = `${baseName}ResourceModel`;

  return `
// Code generated by providergen. DO NOT EDIT.
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
  resp.TypeName = req.ProviderTypeName + "_${resource.name}"
}

func (r *${resourceName}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: ${JSON.stringify(resource.description ?? "")},
    ${generateTerraformSchemaAttributesBlocks({
      mode: resource.mode,
      name: `${camelize(resource.name)}ResourceModel`,
      attributes: [resource.idAttribute, ...resource.attributes],
    })}
  }

  if ext, ok := any(r).(modifySchemaResponseExtension); ok {
    ext.modifySchemaResponse(ctx, resp)
  }
}

func (r *${resourceName}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
  var data ${modelName}

  // Read Terraform plan data into the model
  resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }

  // Create API call logic
  modelIn := tfutils.MergeDiagnostics(data.ToApi(ctx))(&resp.Diagnostics)
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

  resp.Diagnostics.Append(data.FromApi(ctx, modelOut)...)
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
    data.Id.ValueString(),${resource.getHasQueryParams ? "\nnil," : ""}
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

  resp.Diagnostics.Append(data.FromApi(ctx, modelOut)...)
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

  ${generateResourceUpdateDiff(resource)}

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

  resp.Diagnostics.Append(data.FromApi(ctx, modelOut)...)
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

${generateResourceModel(resource)}
`;
}

function generatePrimitiveTypeV2({
  parent,
  attribute,
}: {
  parent: string;
  attribute: IRType;
}) {
  return match(attribute)
    .with({ type: "string", nullable: true }, () => "*string")
    .with({ type: "string" }, () => "string")
    .with({ type: "int", nullable: true }, () => "*int64")
    .with({ type: "int" }, () => "int64")
    .with({ type: "bool", nullable: true }, () => "*bool")
    .with({ type: "bool" }, () => "bool")
    .with(
      { type: "set", elementType: "string", nullable: true },
      () => "*[]string"
    )
    .with({ type: "set", elementType: "string" }, () => "[]string")
    .with(
      { type: "object", nullable: true },
      () => `*${parent}${camelize(attribute.name)}`
    )
    .with({ type: "object" }, () => `${parent}${camelize(attribute.name)}`)
    .exhaustive();
}

function generateClientModel({
  resourceType,
  name,
  attributes,
}: {
  resourceType?: string;
  name: string;
  attributes: IRType[];
}) {
  const lines: string[] = [];
  const extras: string[] = [];

  if (resourceType) {
    lines.push(
      `Id string \`jsonapi:"primary,${pluralize(resourceType)},omitempty"\``
    );
  }

  for (const attribute of attributes) {
    lines.push(
      `${camelize(attribute.name)} ${generatePrimitiveTypeV2({ parent: name, attribute })} \`jsonapi:"attribute" json:"${underscore(attribute.name)},omitempty"\``
    );

    extras.push(
      ...match(attribute)
        .with({ type: "object" }, (attribute) => [
          generateClientModel({
            name: `${name}${camelize(attribute.name)}`,
            attributes: attribute.attributes,
          }),
        ])
        .otherwise(() => [])
    );
  }

  return dedent`
type ${name} struct {
  ${lines.join("\n")}
}

${extras.join("\n\n")}
`;
}

export function generateResourceClientModel(resource: ResourceIR) {
  const name = `${camelize(resource.name)}ResourceModel`;

  return `
// Code generated by providergen. DO NOT EDIT.
package apiclient

${generateClientModel({ resourceType: resource.name, name, attributes: resource.attributes })}
`;
}

export function generateProvider({ resources }: { resources: ResourceIR[] }) {
  return `
// Code generated by providergen. DO NOT EDIT.
package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

var RootlyProviderGeneratedResources = []func() resource.Resource{
  ${resources
    .map((resource) => `New${camelize(resource.name)}Resource,`)
    .join("\n")}
}
`;
}
