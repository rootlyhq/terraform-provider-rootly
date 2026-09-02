import { camelize } from "inflection";
import type { oas30 } from "openapi3-ts";
import { match, P } from "ts-pattern";

export function tfAttributeSchemaType({
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
      {
        type: "array",
        items: { type: P.union("string", "integer") },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) => `schema.${camelize(schema["x-tf-collection-type"])}Attribute`,
    )
    .with(
      {
        type: "array",
        items: { type: "object" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `schema.${camelize(schema["x-tf-collection-type"])}NestedAttribute`,
    )
    .exhaustive();
}

export function tfAttributeValidatorType({
  schema,
}: {
  schema: oas30.SchemaObject;
}) {
  return match(schema)
    .with({ type: "string" }, () => "validator.String")
    .with({ type: "boolean" }, () => "validator.Bool")
    .with({ type: "integer" }, () => "validator.Int64")
    .with({ type: "object" }, () => "validator.Object")
    .with({ type: "array" }, () => "validator.Set")
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
      {
        type: "array",
        items: { type: "string" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.${camelize(schema["x-tf-collection-type"])}ValueOf[string]`,
    )
    .with(
      {
        type: "array",
        items: { type: "integer" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.${camelize(schema["x-tf-collection-type"])}ValueOf[int64]`,
    )
    .with(
      {
        type: "array",
        items: { type: "object" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.${camelize(schema["x-tf-collection-type"])}NestedObjectValueOf[${parent}${name}Item]`,
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
      {
        type: "array",
        items: { type: "string" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.New${camelize(schema["x-tf-collection-type"])}TypeOf[string](ctx)`,
    )
    .with(
      {
        type: "array",
        items: { type: "integer" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.New${camelize(schema["x-tf-collection-type"])}TypeOf[int64](ctx)`,
    )
    .with(
      {
        type: "array",
        items: { type: "object" },
        "x-tf-collection-type": P.union("list", "set"),
      },
      (schema) =>
        `supertypes.New${camelize(schema["x-tf-collection-type"])}NestedObjectTypeOf[${parent}${camelize(name)}Item](ctx)`,
    )
    .otherwise(() => null);
}

export function tfAttributeDefault({ schema }: { schema: oas30.SchemaObject }) {
  if (!schema.default) {
    return null;
  }

  return match(schema)
    .with(
      { type: "string", default: P.string },
      (schema) =>
        `stringdefault.StaticString(${JSON.stringify(schema.default)})`,
    )
    .exhaustive();
}
