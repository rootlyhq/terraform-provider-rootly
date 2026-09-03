import { camelize } from "inflection";
import { match, P } from "ts-pattern";
import type { AttributeType } from "./schema";

export function tfAttributeSchemaType({
  attribute,
  type,
}: {
  attribute: AttributeType;
  type: "attribute" | "block";
}) {
  return match([attribute, type])
    .with([{ type: "string" }, P.any], () => "schema.StringAttribute")
    .with([{ type: "bool" }, P.any], () => "schema.BoolAttribute")
    .with([{ type: "int64" }, P.any], () => "schema.Int64Attribute")
    .with([{ type: "object" }, P.any], () => "schema.SingleNestedAttribute")
    .with([{ type: "list" }, P.any], () => "schema.ListAttribute")
    .with(
      [{ type: "list_nested" }, "attribute"],
      () => "schema.ListNestedAttribute",
    )
    .with([{ type: "list_nested" }, "block"], () => "schema.ListNestedBlock")
    .with([{ type: "set" }, P.any], () => "schema.SetAttribute")
    .with(
      [{ type: "set_nested" }, "attribute"],
      () => "schema.SetNestedAttribute",
    )
    .with([{ type: "set_nested" }, "block"], () => "schema.SetNestedBlock")
    .exhaustive();
}

export function tfAttributeValidatorType({
  attribute,
}: {
  attribute: AttributeType;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "validator.String")
    .with({ type: "bool" }, () => "validator.Bool")
    .with({ type: "int64" }, () => "validator.Int64")
    .with({ type: "object" }, () => "validator.Object")
    .with({ type: "list" }, () => "validator.List")
    .with({ type: "list_nested" }, () => "validator.List")
    .with({ type: "set" }, () => "validator.Set")
    .with({ type: "set_nested" }, () => "validator.Set")
    .exhaustive();
}

export function tfAttributePlanModifierType({
  attribute,
}: {
  attribute: AttributeType;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "planmodifier.String")
    .with({ type: "int64" }, () => "planmodifier.Int64")
    .with({ type: "bool" }, () => "planmodifier.Bool")
    .exhaustive();
}

export function tfAttributeValueType({
  attribute,
  parent,
}: {
  attribute: AttributeType;
  parent: string;
}) {
  return match(attribute)
    .with({ type: "string" }, () => "types.String")
    .with({ type: "bool" }, () => "types.Bool")
    .with({ type: "int64" }, () => "types.Int64")
    .with(
      { type: "list" },
      (attribute) => `supertypes.ListValueOf[${attribute.elementType}]`,
    )
    .with(
      { type: "set" },
      (attribute) => `supertypes.SetValueOf[${attribute.elementType}]`,
    )
    .with(
      { type: "set_nested" },
      (attribute) =>
        `supertypes.SetNestedObjectValueOf[${parent}${camelize(attribute.name)}Item]`,
    )
    .with(
      { type: "list_nested" },
      (attribute) =>
        `supertypes.ListNestedObjectValueOf[${parent}${camelize(attribute.name)}Item]`,
    )
    .exhaustive();
}

export function tfAttributeCustomType({
  parent,
  attribute,
}: {
  parent: string;
  attribute: AttributeType;
}) {
  return match(attribute)
    .with(
      { type: "object" },
      (value) =>
        `supertypes.NewSingleNestedObjectTypeOf[${parent}${camelize(value.name)}](ctx)`,
    )
    .with(
      { type: "list", elementType: P.string },
      (value) => `supertypes.NewListTypeOf[${value.elementType}](ctx)`,
    )
    .with(
      { type: "set", elementType: "string" },
      (value) => `supertypes.NewSetTypeOf[${value.elementType}](ctx)`,
    )
    .with(
      { type: "list_nested" },
      (value) =>
        `supertypes.NewListNestedObjectTypeOf[${parent}${camelize(value.name)}Item](ctx)`,
    )
    .with(
      { type: "set_nested" },
      (value) =>
        `supertypes.NewSetNestedObjectTypeOf[${parent}${camelize(value.name)}Item](ctx)`,
    )
    .otherwise(() => null);
}

export function tfAttributeDefault({
  attribute,
}: {
  attribute: AttributeType;
}) {
  if (!("default" in attribute) || !attribute.default) {
    return null;
  }

  return match(attribute)
    .with(
      { type: "string", default: P.string },
      (schema) =>
        `stringdefault.StaticString(${JSON.stringify(schema.default)})`,
    )
    .exhaustive();
}
