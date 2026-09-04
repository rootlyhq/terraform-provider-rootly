import {
  buildValidators,
  resolveDescription,
  type AttributeBlockType,
  type AttributeType,
  type DataSourceDef,
  type ResourceDef,
} from "./schema";
import {
  tfAttributeCustomType,
  tfAttributeDefault,
  tfAttributePlanModifierType,
  tfAttributeSchemaType,
  tfAttributeValidatorType,
  tfAttributeValueType,
} from "./go-types";
import { camelize } from "inflection";
import { match, P } from "ts-pattern";

export function generateSchemaAttributes({
  parent,
  attributes,
  blocks,
}: {
  parent: string;
  attributes: AttributeType[];
  blocks: AttributeBlockType[];
}) {
  const lines: string[] = [];

  if (attributes.length > 0) {
    lines.push("Attributes: map[string]schema.Attribute{");
    for (const attribute of attributes) {
      lines.push(
        `"${attribute.name}": ${generateSchemaAttribute({ parent, attribute, type: "attribute" })},`,
      );
    }
    lines.push("},");
  }

  if (blocks.length > 0) {
    lines.push("Blocks: map[string]schema.Block{");
    for (const block of blocks) {
      lines.push(
        `"${block.name}": ${generateSchemaAttribute({ parent, attribute: block, type: "block" })},`,
      );
    }
    lines.push("},");
  }

  return lines.join("\n");
}

function generateSchemaAttribute({
  parent,
  attribute,
  type,
}: {
  parent: string;
  attribute: AttributeType;
  type: "attribute" | "block";
}) {
  const parts: string[] = [];
  parts.push(`${tfAttributeSchemaType({ attribute, type })}{`);

  const description = resolveDescription(attribute);
  if (description) {
    parts.push(`MarkdownDescription: ${JSON.stringify(description)},`);
  }

  if (type === "attribute") {
    parts.push(
      ...match(attribute.computedOptionalRequired)
        .with("required", () => ["Required: true,"])
        .with("computed", () => ["Computed: true,"])
        .with("computed_optional", () => ["Optional: true,", "Computed: true,"])
        .with("optional", () => ["Optional: true,"])
        .exhaustive(),
    );
  }

  const tfCustomType = tfAttributeCustomType({ parent, attribute });
  if (tfCustomType) {
    parts.push(`CustomType: ${tfCustomType},`);
  }

  const tfDefault = tfAttributeDefault({ attribute });
  if (tfDefault) {
    parts.push(`Default: ${tfDefault},`);
  }

  const validators = buildValidators(attribute);
  if (validators.length > 0) {
    const tfValidatorType = tfAttributeValidatorType({ attribute });
    parts.push(`Validators: []${tfValidatorType}{`);
    parts.push(...validators.map((v) => `${v},`));
    parts.push(`},`);
  }

  if (attribute.planModifiers && attribute.planModifiers.length > 0) {
    const tfPlanModifierType = tfAttributePlanModifierType({ attribute });
    parts.push(`PlanModifiers: []${tfPlanModifierType}{`);
    parts.push(...attribute.planModifiers.map((v) => `${v},`));
    parts.push(`},`);
  }

  match([attribute, type])
    .with([{ type: "single_nested" }, P.any], ([value]) => {
      parts.push(
        generateSchemaAttributes({
          parent: `${parent}${camelize(value.name)}`,
          attributes: value.attributes,
          blocks: value.blocks,
        }),
      );
    })
    .with(
      [{ type: "list_nested" }, "attribute"],
      [{ type: "set_nested" }, "attribute"],
      ([value]) => {
        parts.push("NestedObject: schema.NestedAttributeObject{");
        parts.push(
          generateSchemaAttributes({
            parent: `${parent}${camelize(value.name)}Item`,
            attributes: value.attributes,
            blocks: value.blocks,
          }),
        );
        parts.push("},");
      },
    )
    .with(
      [{ type: "list_nested" }, "block"],
      [{ type: "set_nested" }, "block"],
      ([value]) => {
        parts.push("NestedObject: schema.NestedBlockObject{");
        parts.push(
          generateSchemaAttributes({
            parent: `${parent}${camelize(value.name)}Item`,
            attributes: value.attributes,
            blocks: value.blocks,
          }),
        );
        parts.push("},");
      },
    );

  parts.push(`}`);

  return parts.join("\n");
}

export function generateModels({
  def,
  name,
  clientName,
  attributes,
  level,
  options,
}: {
  def: DataSourceDef | ResourceDef;
  name: string;
  clientName: string;
  attributes: AttributeType[];
  level: number;
  options: {
    rootIsCollection: boolean;
    generateFromApi: boolean;
    generateToApiForCreate: boolean;
    generateToApiForUpdate: boolean;
  };
}) {
  const fields: string[] = [];
  const children: string[] = [];

  for (const attribute of attributes) {
    const goAttributeValueType = tfAttributeValueType({
      attribute,
      parent: name,
    });

    fields.push(
      `${camelize(attribute.name)} ${goAttributeValueType} \`tfsdk:"${attribute.name}"\``,
    );

    match(attribute)
      .with({ type: "single_nested" }, (attribute) => {
        children.push(
          generateModels({
            def,
            name: `${name}${camelize(attribute.name)}`,
            clientName: `${clientName}${camelize(attribute.name)}`,
            attributes: [...attribute.attributes, ...attribute.blocks],
            level: level + 1,
            options,
          }),
        );
      })
      .with({ type: "list_nested" }, { type: "set_nested" }, (attribute) => {
        children.push(
          generateModels({
            def,
            name: `${name}${camelize(attribute.name)}Item`,
            clientName: attribute.hints?.isTopLevelCollection
              ? clientName
              : `${clientName}${camelize(attribute.name)}Item`,
            attributes: [...attribute.attributes, ...attribute.blocks],
            level: level + 1,
            options,
          }),
        );
      });
  }

  return `
type ${name} struct {
${fields.join("\n")}
}

${
  options.generateFromApi
    ? generateModelFromApi({
        name,
        clientName,
        attributes,
        level,
        rootIsCollection: options.rootIsCollection,
      })
    : ""
}

${
  options.generateToApiForCreate
    ? generateModelToApi({
        action: "create",
        name,
        clientName,
        attributes,
      })
    : ""
}

${
  options.generateToApiForUpdate
    ? generateModelToApi({
        action: "update",
        name,
        clientName,
        attributes,
      })
    : ""
}

${children.join("\n\n")}
`;
}

function generateModelFromApi({
  name,
  clientName,
  attributes,
  level,
  rootIsCollection,
}: {
  name: string;
  clientName: string;
  attributes: AttributeType[];
  level: number;
  rootIsCollection: boolean;
}) {
  const dataType =
    rootIsCollection && level === 0
      ? `[]apiclient.${clientName}`
      : `apiclient.${clientName}`;

  const assignments = attributes.map((attribute) => {
    if (
      !attribute.schemas?.read &&
      !(attribute.type === "string" && attribute.hints?.isOpenApiId)
    ) {
      return `// ${attribute.name} is not returned`;
    }
    const goValue = match(attribute)
      .with(
        { type: "string", hints: { isOpenApiId: true } },
        (attribute) => `types.StringValue(data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "string" },
        (attribute) =>
          `jsonapitypes.NullableStringValue(data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "int64" },
        (attribute) =>
          `jsonapitypes.NullableInt64Value(data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "bool" },
        (attribute) =>
          `jsonapitypes.NullableBoolValue(data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "list" },
        (attribute) =>
          `jsonapitypes.NullableListValueOfSlice(ctx, data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "set" },
        (attribute) =>
          `jsonapitypes.NullableSetValueOfSlice(ctx, data.${camelize(attribute.name)})`,
      )
      .with(
        { type: "list_nested" },
        (attribute) =>
          `diagutils.MergeDiagnostics(jsonapitypes.ConvertToListModel(
  ctx,
  ${attribute.hints?.isTopLevelCollection ? "jsonapi.NewNullableAttrWithValue(data)" : `data.${camelize(attribute.name)}`},
  func(ctx context.Context, item *${name}${camelize(attribute.name)}Item, apiItem ${attribute.hints?.isTopLevelCollection ? `apiclient.${clientName}` : `apiclient.${clientName}${camelize(attribute.name)}Item`}) diag.Diagnostics {
    return item.FromApi(ctx, apiItem)
  },
))(&diags)`,
      )
      .with(
        { type: "set_nested" },
        (attribute) =>
          `diagutils.MergeDiagnostics(jsonapitypes.ConvertToSetModel(
  ctx,
  data.${camelize(attribute.name)},
  func(ctx context.Context, item *${name}${camelize(attribute.name)}Item, apiItem apiclient.${clientName}${camelize(attribute.name)}Item) diag.Diagnostics {
    return item.FromApi(ctx, apiItem)
  },
))(&diags)`,
      )
      .with(
        { type: "single_nested" },
        (attribute) =>
          `diagutils.MergeDiagnostics(jsonapitypes.ConvertToSingleModel(
  ctx,
  data.${camelize(attribute.name)},
  func(ctx context.Context, item *${name}${camelize(attribute.name)}, apiItem apiclient.${clientName}${camelize(attribute.name)}) diag.Diagnostics {
    return item.FromApi(ctx, apiItem)
  },
))(&diags)`,
      )
      .otherwise(
        (attribute) => `nil // TODO: Implement: ${JSON.stringify(attribute)}`,
      );

    return `m.${camelize(attribute.name)} = ${goValue}`;
  });

  return `
func (m *${name}) FromApi(ctx context.Context, data ${dataType}) diag.Diagnostics {
	var diags diag.Diagnostics

	${assignments.join("\n")}

	return diags
}
`;
}

function generateModelToApi({
  action,
  name,
  clientName,
  attributes,
}: {
  action: "create" | "update";
  name: string;
  clientName: string;
  attributes: AttributeType[];
}) {
  const assignments = attributes.map((attribute) => {
    if (!attribute.schemas?.[action]) {
      return `// ${attribute.name} is not available for ${action}`;
    }
    const goValue = match(attribute)
      .with(
        { type: "string" },
        (attribute) =>
          `jsonapitypes.NewNullableFromString(m.${camelize(attribute.name)})`,
      )
      .with(
        { type: "int64" },
        (attribute) =>
          `jsonapitypes.NewNullableFromInt64(m.${camelize(attribute.name)})`,
      )
      .with(
        { type: "bool" },
        (attribute) =>
          `jsonapitypes.NewNullableFromBool(m.${camelize(attribute.name)})`,
      )
      .with(
        { type: "list" },
        (attribute) =>
          `diagutils.MergeDiagnostics(jsonapitypes.NewNullableFromListOf(ctx, m.${camelize(attribute.name)}))(&diags)`,
      )
      .with(
        { type: "set" },
        (attribute) =>
          `diagutils.MergeDiagnostics(jsonapitypes.NewNullableFromSetOf(ctx, m.${camelize(attribute.name)}))(&diags)`,
      )
      .with({ type: "list_nested" }, { type: "set_nested" }, (attribute) => {
        const unknownBehavior = attribute.hacks?.unknownBehavior ?? "omit";
        const nullBehavior = attribute.hacks?.nullBehavior ?? "null";
        const emptyBehavior = attribute.hacks?.emptyBehavior ?? "empty";
        const resolveOutcome = (behavior: "omit" | "null" | "empty"): string =>
          match(behavior)
            .with("omit", () => "jsonapitypes.OutcomeOmit")
            .with("null", () => "jsonapitypes.OutcomeNull")
            .with("empty", () => "jsonapitypes.OutcomeEmptyList")
            .exhaustive();

        return `diagutils.MergeDiagnostics(jsonapitypes.ConvertNullableList(
  ctx,
  m.${camelize(attribute.name)},
  jsonapitypes.NullableListConfig{
    OnUnknown: ${resolveOutcome(unknownBehavior)},
    OnNull: ${resolveOutcome(nullBehavior)},
    OnEmpty: ${resolveOutcome(emptyBehavior)},
  },
  func(ctx context.Context, item *${name}${camelize(attribute.name)}Item) (*apiclient.${clientName}${camelize(attribute.name)}Item, diag.Diagnostics) {
    return item.ToApiFor${camelize(action)}(ctx)
  },
))(&diags)`;
      })
      .otherwise(
        (attribute) => `nil // TODO: Implement: ${JSON.stringify(attribute)}`,
      );

    return `data.${camelize(attribute.name)} = ${goValue}`;
  });

  return `
func (m *${name}) ToApiFor${camelize(action)}(ctx context.Context) (*apiclient.${clientName}, diag.Diagnostics) {
	var diags diag.Diagnostics

  data := apiclient.${clientName}{}
	${assignments.join("\n")}

	return &data, diags
}
`;
}
