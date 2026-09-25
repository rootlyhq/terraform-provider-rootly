import { camelize } from "inflection";
import type { oas30 } from "openapi3-ts";
import { match } from "ts-pattern";

export function tfSchemaAttributeType({
  schema,
}: {
  schema: oas30.SchemaObject;
}) {
  return match(schema)
    .with({ type: "string" }, () => "schema.StringAttribute")
    .with({ type: "boolean" }, () => "schema.BoolAttribute")
    .with({ type: "integer" }, () => "schema.Int64Attribute")
    .with({ type: "object" }, () => "schema.SingleNestedAttribute")
    .with(
      { type: "array", items: { type: "string" } },
      { type: "array", items: { type: "integer" } },
      () => "schema.ListAttribute",
    )
    .with(
      { type: "array", items: { type: "object" } },
      () => "schema.ListNestedAttribute",
    )
    .exhaustive();
}

export function tfAttributeValueType({
  schema,
  parent,
  name,
}: {
  schema: oas30.SchemaObject;
  parent: string;
  name: string;
}) {
  return match(schema)
    .with({ type: "string" }, () => "types.String")
    .with({ type: "boolean" }, () => "types.Bool")
    .with({ type: "integer" }, () => "types.Int64")
    .with(
      { type: "object" },
      () => `supertypes.SingleNestedObjectValueOf[${parent}${name}]`,
    )
    .with(
      { type: "array", items: { type: "string" } },
      () => "supertypes.ListValueOf[string]",
    )
    .with(
      { type: "array", items: { type: "integer" } },
      () => "supertypes.ListValueOf[int64]",
    )
    .with(
      { type: "array", items: { type: "object" } },
      () => `supertypes.ListNestedObjectValueOf[${parent}${name}Item]`,
    )
    .exhaustive();
}

export function tfAttributeCustomType({
  schema,
  parent,
  name,
}: {
  schema: oas30.SchemaObject;
  parent: string | null;
  name: string;
}) {
  return match(schema)
    .with(
      { type: "object" },
      () =>
        `supertypes.NewSingleNestedObjectTypeOf[${parent}${camelize(name)}](ctx)`,
    )
    .with(
      { type: "array", items: { type: "string" } },
      () => "supertypes.NewListTypeOf[string](ctx)",
    )
    .with(
      { type: "array", items: { type: "integer" } },
      () => "supertypes.NewListTypeOf[int64](ctx)",
    )
    .with(
      { type: "array", items: { type: "object" } },
      () =>
        `supertypes.NewListNestedObjectTypeOf[${parent}${camelize(name)}Item](ctx)`,
    )
    .otherwise(() => null);
}
