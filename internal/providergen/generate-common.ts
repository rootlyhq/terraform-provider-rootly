import {
  isCollectionAttribute,
  type AttributeType,
  type DataSourceDef,
  type ResourceDef,
} from "./schema";
import { tfAttributeValueType } from "./go-types";
import { camelize } from "inflection";
import { match, P } from "ts-pattern";

export function generateModels({
  def,
  name,
  clientName,
  attributes,
}: {
  def: DataSourceDef | ResourceDef;
  name: string;
  clientName: string;
  attributes: AttributeType[];
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

    if (isCollectionAttribute(attribute)) {
      children.push(
        generateModels({
          def,
          name: `${name}${camelize(attribute.name)}Item`,
          clientName: `${clientName}${camelize(attribute.name)}Item`,
          attributes: [...attribute.attributes, ...attribute.blocks],
        }),
      );
    }
  }

  return `
type ${name} struct {
${fields.join("\n")}
}

${generateFromApi({
  name,
  clientName,
  attributes,
})}

${generateToApi({
  action: "create",
  name,
  clientName,
  attributes,
})}

${generateToApi({
  action: "update",
  name,
  clientName,
  attributes,
})}

${children.join("\n\n")}
`;
}

function generateFromApi({
  name,
  clientName,
  attributes,
}: {
  name: string;
  clientName: string;
  attributes: AttributeType[];
}) {
  const assignments = attributes.map((attribute) => {
    if (!attribute.schemas?.read) {
      return `// ${attribute.name} is not returned`;
    }
    const goValue = match(attribute)
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
        (
          attribute,
        ) => `(func() supertypes.ListNestedObjectValueOf[${name}${camelize(attribute.name)}Item] {
  if v, err := data.${camelize(attribute.name)}.Get(); err == nil {
    return supertypes.NewListNestedObjectValueOfValueSlice(ctx, lo.Map(v, func(vv apiclient.${clientName}${camelize(attribute.name)}Item, _ int) ${name}${camelize(attribute.name)}Item {
      var mm ${name}${camelize(attribute.name)}Item
      diags.Append(mm.FromApi(ctx, vv)...)
      return mm
    }))
  }
  return supertypes.NewListNestedObjectValueOfNull[${name}${camelize(attribute.name)}Item](ctx)
})()`,
      )
      .with(
        { type: "set_nested" },
        (
          attribute,
        ) => `(func() supertypes.SetNestedObjectValueOf[${name}${camelize(attribute.name)}Item] {
  if v, err := data.${camelize(attribute.name)}.Get(); err == nil {
    return supertypes.NewSetNestedObjectValueOfValueSlice(ctx, lo.Map(v, func(vv apiclient.${clientName}${camelize(attribute.name)}Item, _ int) ${name}${camelize(attribute.name)}Item {
      var mm ${name}${camelize(attribute.name)}Item
      diags.Append(mm.FromApi(ctx, vv)...)
      return mm
    }))
  }
  return supertypes.NewSetNestedObjectValueOfNull[${name}${camelize(attribute.name)}Item](ctx)
})()`,
      )
      .otherwise(
        (attribute) => `nil // TODO: Implement: ${JSON.stringify(attribute)}`,
      );

    return `m.${camelize(attribute.name)} = ${goValue}`;
  });

  return `
func (m *${name}) FromApi(ctx context.Context, data apiclient.${clientName}) diag.Diagnostics {
	var diags diag.Diagnostics

	${assignments.join("\n")}

	return diags
}
`;
}

function generateToApi({
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
