import type { oas30 } from "openapi3-ts";
import type {
  ResolvedDataSourceConfig,
  ResolvedResourceConfig,
} from "./schema";
import { assertSchemaObject } from "./types";
import { tfAttributeValueType } from "./go-types";
import { camelize } from "inflection";
import { match, P } from "ts-pattern";

export function generateModel({
  config,
  schema,
  baseName,
  name,
  isTopLevel,
}: {
  config: ResolvedDataSourceConfig | ResolvedResourceConfig;
  schema: oas30.SchemaObject;
  baseName: string;
  name: string;
  isTopLevel: boolean;
}) {
  const fromApiDataType =
    config.kind === "data_source" && config.type === "list" && isTopLevel
      ? `[]apiclient.${baseName}`
      : `apiclient.${baseName}`;
  if (!schema.properties) {
    throw new Error(`Schema ${name} does not have properties`);
  }

  const required = new Set(schema.required);

  const fields: string[] = [];
  const fromApiAssignments: string[] = [];
  const children: string[] = [];

  for (const [key, value] of Object.entries(schema.properties)) {
    assertSchemaObject(value);
    const property = {
      required: required.has(key),
      nullable: value.nullable ?? false,
    };

    // Terraform Attribute
    const goType = tfAttributeValueType({
      schema: value,
      parent: name,
      name: camelize(key),
    });

    if (goType) {
      fields.push(`${camelize(key)} ${goType} \`tfsdk:"${key}"\``);
    } else {
      fields.push(
        `// TODO: Support property ${key} of type ${JSON.stringify(value)}`,
      );
    }

    match(value)
      .with({ type: "object" }, (value) => {
        children.push(
          generateModel({
            config,
            schema: value,
            baseName: `${baseName}${camelize(key)}`,
            name: `${name}${camelize(key)}`,
            isTopLevel: false,
          }),
        );
      })
      .with(
        {
          type: "array",
          items: { type: "object" },
          "x-tf-top-level-item-type": true,
        },
        (value) => {
          children.push(
            generateModel({
              config,
              schema: value.items as oas30.SchemaObject,
              baseName,
              name: `${name}${camelize(key)}Item`,
              isTopLevel: false,
            }),
          );
        },
      )
      .with({ type: "array", items: { type: "object" } }, (value) => {
        children.push(
          generateModel({
            config,
            schema: value.items as oas30.SchemaObject,
            baseName: `${baseName}${camelize(key)}Item`,
            name: `${name}${camelize(key)}Item`,
            isTopLevel: false,
          }),
        );
      });

    // FromApi Assignment
    fromApiAssignments.push(
      match([value, property])
        .with(
          [{ type: "string" }, { nullable: false }],
          () => `m.${camelize(key)} = types.StringValue(data.${camelize(key)})`,
        )
        .with(
          [{ type: "string" }, { nullable: true }],
          () =>
            `m.${camelize(key)} = jsonapitypes.NullableStringValue(data.${camelize(key)})`,
        )
        .with(
          [{ type: "boolean" }, { nullable: false }],
          () => `m.${camelize(key)} = types.BoolValue(data.${camelize(key)})`,
        )
        .with(
          [{ type: "boolean" }, { nullable: true }],
          () =>
            `m.${camelize(key)} = jsonapitypes.NullableBoolValue(data.${camelize(key)})`,
        )
        .with(
          [{ type: "integer" }, { nullable: false }],
          () => `m.${camelize(key)} = types.Int64Value(data.${camelize(key)})`,
        )
        .with(
          [{ type: "integer" }, { nullable: true }],
          () =>
            `m.${camelize(key)} = jsonapitypes.NullableInt64Value(data.${camelize(key)})`,
        )
        .with(
          [{ type: "object" }, { nullable: true }],
          () =>
            `m.${camelize(key)} = (func() supertypes.SingleNestedObjectValueOf[${name}${camelize(key)}] {
  if v, err := data.${camelize(key)}.Get(); err == nil {
    var mm ${name}${camelize(key)}
    diags.Append(mm.FromApi(ctx, v)...)
		return supertypes.NewSingleNestedObjectValueOf(ctx, &mm)
	}
	return supertypes.NewSingleNestedObjectValueOfNull[${name}${camelize(key)}](ctx)
})()`,
        )
        .with(
          [
            {
              type: "array",
              items: {
                type: P.union("string", "boolean", "integer"),
              },
            },
            { nullable: true },
          ],
          () =>
            `m.${camelize(key)} = jsonapitypes.NullableListValueOfSlice(ctx, data.${camelize(key)})`,
        )
        .with(
          [
            {
              type: "array",
              items: {
                type: P.union("string", "boolean", "integer"),
              },
            },
            { nullable: false },
          ],
          () =>
            `m.${camelize(key)} = supertypes.NewListValueOfSlice(ctx, data.${camelize(key)})`,
        )
        .with(
          [
            {
              type: "array",
              items: { type: "object" },
              "x-tf-top-level-item-type": true,
            },
            { nullable: false },
          ],
          () => `m.${camelize(key)} = (func() supertypes.ListNestedObjectValueOf[${name}${camelize(key)}Item] {
  return supertypes.NewListNestedObjectValueOfValueSlice(ctx, lo.Map(data, func(vv apiclient.${baseName}, _ int) ${name}${camelize(key)}Item {
    var mm ${name}${camelize(key)}Item
    diags.Append(mm.FromApi(ctx, vv)...)
    return mm
  }))
})()`,
        )
        .with(
          [{ type: "array", items: { type: "object" } }, { nullable: false }],
          () => `m.${camelize(key)} = (func() supertypes.ListNestedObjectValueOf[${name}${camelize(key)}Item] {
  return supertypes.NewListNestedObjectValueOfValueSlice(ctx, lo.Map(data.${camelize(key)}, func(vv apiclient.${baseName}${camelize(key)}Item, _ int) ${name}${camelize(key)}Item {
    var mm ${name}${camelize(key)}Item
    diags.Append(mm.FromApi(ctx, vv)...)
    return mm
  }))
})()`,
        )
        .with(
          [{ type: "array", items: { type: "object" } }, { nullable: true }],
          () => `m.${camelize(key)} = (func() supertypes.ListNestedObjectValueOf[${name}${camelize(key)}Item] {
  if v, err := data.${camelize(key)}.Get(); err == nil {
    return supertypes.NewListNestedObjectValueOfValueSlice(ctx, lo.Map(v, func(vv apiclient.${baseName}${camelize(key)}Item, _ int) ${name}${camelize(key)}Item {
      var mm ${name}${camelize(key)}Item
      diags.Append(mm.FromApi(ctx, vv)...)
      return mm
    }))
  }
  return supertypes.NewListNestedObjectValueOfNull[${name}${camelize(key)}Item](ctx)
})()`,
        )
        .otherwise(
          () =>
            `// TODO: Implement FromApi for ${key} of type ${JSON.stringify(value)} and property ${JSON.stringify(property)}`,
        ),
    );
  }

  return `
type ${name} struct {
${fields.join("\n")}
}

func (m *${name}) FromApi(ctx context.Context, data ${fromApiDataType}) diag.Diagnostics {
	var diags diag.Diagnostics

	${fromApiAssignments.join("\n")}

	return diags
}

${children.join("\n\n")}
`;
}
