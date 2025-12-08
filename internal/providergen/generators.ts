import { match } from "ts-pattern";
import type { IRObject, IRResource, IRType, ResourceIR } from "./ir";
import { camelize, pluralize, underscore } from "inflection";
import dedent from "dedent";

function generateTerraformValueType({
  parent,
  attribute,
}: {
  parent: string;
  attribute: IRType;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "supertypes.StringValue")
    .with({ type: "int" }, () => "supertypes.Int64Value")
    .with({ type: "bool" }, () => "supertypes.BoolValue")
    .with(
      { type: "list", elementType: "string" },
      () => "supertypes.ListValueOf[string]"
    )
    .with(
      { type: "set", elementType: "string" },
      () => "supertypes.SetValueOf[string]"
    )
    .with(
      { type: "set_nested" },
      () =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(attribute.name)}Item]`
    )
    .with(
      { type: "object" },
      () =>
        `supertypes.SingleNestedObjectValueOf[${parent}${camelize(attribute.name)}]`
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

function generateTerraformAttribute({
  parent,
  attribute,
}: {
  parent: string;
  attribute: IRType;
}) {
  let description = attribute.description ?? "";
  if (attribute.deprecationMessage) {
    description += ` **Deprecated** ${attribute.deprecationMessage}`;
  }

  const commonParts: string[] = [];
  commonParts.push(`MarkdownDescription: ${JSON.stringify(description)},`);
  if (attribute.deprecationMessage) {
    commonParts.push(
      `DeprecationMessage: ${JSON.stringify(attribute.deprecationMessage)},`
    );
  }
  commonParts.push(
    match(attribute.computedOptionalRequired)
      .with("required", () => "Required: true,")
      .with("computed", () => "Computed: true,")
      .with("computed_optional", () => "Optional: true,\nComputed: true,")
      .with("optional", () => "Optional: true,")
      .exhaustive()
  );
  if (attribute.sensitive) {
    commonParts.push("Sensitive: true,");
  }

  return match(attribute)
    .with({ type: "string" }, () => {
      const parts: string[] = [];
      parts.push("schema.StringAttribute{");
      parts.push(...commonParts);
      parts.push("CustomType: supertypes.StringType{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.String{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
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
      parts.push("CustomType: supertypes.Int64Type{},");
      if (attribute.validators) {
        parts.push("Validators: []validator.Int64{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
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
      parts.push("CustomType: supertypes.NewListTypeOf[string](ctx),");
      if (attribute.validators) {
        parts.push("Validators: []validator.List{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
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
      parts.push("CustomType: supertypes.NewSetTypeOf[string](ctx),");
      if (attribute.validators) {
        parts.push("Validators: []validator.Set{");
        parts.push(...attribute.validators.map((validator) => `${validator},`));
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
      parts.push(
        `CustomType: supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(
          attribute.name
        )}Item](ctx),`
      );
      parts.push("NestedObject: schema.NestedAttributeObject{");
      parts.push("Attributes: map[string]schema.Attribute{");
      for (const nestedAttribute of attribute.attributes) {
        parts.push(
          `"${nestedAttribute.name}": ${generateTerraformAttribute({
            parent: `${parent}${camelize(attribute.name)}Item`,
            attribute: nestedAttribute,
          })},`
        );
      }
      parts.push("},");
      parts.push("},");
      parts.push("}");
      return parts.join("\n");
    })
    .with({ type: "object" }, (attribute) => {
      const parts: string[] = [];

      parts.push("schema.SingleNestedAttribute{");
      parts.push(...commonParts);
      parts.push(
        `CustomType: supertypes.NewSingleNestedObjectTypeOf[${parent}${camelize(
          attribute.name
        )}](ctx),`
      );
      parts.push("Attributes: map[string]schema.Attribute{");
      for (const nestedAttribute of attribute.attributes) {
        parts.push(
          `"${nestedAttribute.name}": ${generateTerraformAttribute({
            parent: `${parent}${camelize(attribute.name)}`,
            attribute: nestedAttribute,
          })},`
        );
      }
      parts.push("},");
      parts.push("}");

      return parts.join("\n");
    })
    .exhaustive();
}

function generateResourceSchemaAttributes(resource: ResourceIR) {
  const lines: string[] = [];

  for (const attribute of [resource.idAttribute, ...resource.attributes]) {
    lines.push(
      `"${attribute.name}": ${generateTerraformAttribute({
        parent: `${camelize(resource.name)}ResourceModel`,
        attribute,
      })},`
    );
  }

  return lines.join("\n");
}

function generatePrimitiveToTerraform({
  parent,
  attribute,
  srcVar,
  destVar,
}: {
  parent: string;
  attribute: IRType;
  srcVar: string;
  destVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  const destVarName = `${destVar}.${camelize(attribute.name)}`;
  return match(attribute)
    .with(
      { type: "string", nullable: true },
      () => `${destVarName} = supertypes.NewStringPointerValue(${srcVarName})`
    )
    .with(
      { type: "string", sourceType: "time" },
      () => `${destVarName} = supertypes.NewStringValue(${srcVarName}.String())`
    )
    .with(
      { type: "string" },
      () => `${destVarName} = supertypes.NewStringValue(${srcVarName})`
    )
    .with(
      { type: "int" },
      () => `${destVarName} = supertypes.NewInt64Value(${srcVarName})`
    )
    .with(
      { type: "bool" },
      () => `${destVarName} = supertypes.NewBoolValue(${srcVarName})`
    )
    .with(
      { type: "list", elementType: "string" },
      () =>
        `${destVarName} = supertypes.NewListValueOfSlice(ctx, ${srcVarName})`
    )
    .with(
      { type: "set", nullable: true, elementType: "string" },
      () => dedent`
        if ${srcVarName} == nil {
          ${destVarName} = supertypes.NewSetValueOfNull[string](ctx)
        } else {
          ${destVarName} = supertypes.NewSetValueOfSlice(ctx, *${srcVarName})
        }
      `
    )
    .with(
      { type: "set", elementType: "string" },
      () => `${destVarName} = supertypes.NewSetValueOfSlice(ctx, ${srcVarName})`
    )
    .with(
      { type: "set_nested" },
      (attribute) =>
        `${destVarName} = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, sliceutils.Map(func(item apiclient.${attribute.model}) ${parent}${camelize(attribute.name)}Item {
          var model ${parent}${camelize(attribute.name)}Item
          diags.Append(model.FromApi(ctx, item)...)
          return model
        }, ${srcVarName}))`
    )
    .with(
      { type: "object" },
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
  parent,
  attribute,
  srcVar,
  destVar,
  diagsVar,
}: {
  parent: string;
  attribute: IRType;
  srcVar: string;
  destVar: string;
  diagsVar: string;
}) {
  const srcVarName = `${srcVar}.${camelize(attribute.name)}`;
  const destVarName = `${destVar}.${camelize(attribute.name)}`;
  return match(attribute)
    .with(
      { type: "string", nullable: true },
      () => `${destVarName} = ${srcVarName}.ValueStringPointer()`
    )
    .with(
      { type: "string" },
      () => `${destVarName} = ${srcVarName}.ValueString()`
    )
    .with(
      { type: "int", nullable: true },
      () => `${destVarName} = ${srcVarName}.ValueInt64Pointer()`
    )
    .with({ type: "int" }, () => `${destVarName} = ${srcVarName}.ValueInt64()`)
    .with(
      { type: "bool", nullable: true },
      () => `${destVarName} = ${srcVarName}.ValueBoolPointer()`
    )
    .with({ type: "bool" }, () => `${destVarName} = ${srcVarName}.ValueBool()`)
    .with(
      { type: "list", nullable: true, elementType: "string" },
      () =>
        `${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})`
    )
    .with(
      { type: "list", elementType: "string" },
      () =>
        `${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})`
    )
    .with(
      { type: "set", nullable: true, elementType: "string" },
      () => dedent`
        if ${srcVarName}.IsKnown() {
          ${destVarName} = ptr.Ptr(tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar}))
        } else {
          ${destVarName} = nil
        }
      `
    )
    .with(
      { type: "set", elementType: "string" },
      () =>
        `${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})`
    )
    .with(
      { type: "set_nested" },
      () =>
        `${destVarName} = tfutils.MergeDiagnostics(${srcVarName}.Get(ctx))(&${diagsVar})`
    )
    .with(
      { type: "object", nullable: true },
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
    .otherwise(() => "// TODO");
}

function generateModel({
  name,
  attributes,
}: {
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
        parent: name,
        attribute,
      })} \`tfsdk:"${attribute.name}"\``
    );

    fromApiLines.push(
      `${generatePrimitiveToTerraform({
        parent: name,
        attribute,
        srcVar: "data",
        destVar: "m",
      })}${attribute.deprecationMessage ? " // Deprecated" : ""}`
    );

    toApiLines.push(
      `${generateTerraformToPrimitive({
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
    Attributes: map[string]schema.Attribute{
      ${generateResourceSchemaAttributes(resource)}
    },
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
