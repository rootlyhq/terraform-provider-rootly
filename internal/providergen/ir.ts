import { match, P } from "ts-pattern";

interface IRBase {
  kind: string;
  description?: string;
}

interface IRBaseWithComputedOptionalRequired extends IRBase {
  computedOptionalRequired: "computed" | "optional" | "required";
}

export type IRType =
  | IRString
  | IRBool
  | IRInt
  | IRObject
  | IRArray
  | IRResource;

export interface IRString extends IRBaseWithComputedOptionalRequired {
  kind: "string";
}

export interface IRBool extends IRBaseWithComputedOptionalRequired {
  kind: "bool";
}

export interface IRInt extends IRBaseWithComputedOptionalRequired {
  kind: "int";
}

export interface IRObject extends IRBaseWithComputedOptionalRequired {
  kind: "object";
  fields: Record<string, IRType>;
}

export interface IRArray extends IRBaseWithComputedOptionalRequired {
  kind: "array";
  element: IRType;
}

export interface IRResource extends IRBase {
  kind: "resource";
  resourceType: string;
  idElement: IRType;
  fields: Record<string, IRType>;
}

// TODO: Handle computed
export function toIR({
  schema,
  required,
}: {
  schema: any;
  required: boolean | null;
}): IRType {
  const common = {
    computedOptionalRequired: required ? "required" : "optional",
    description: schema.description,
  } as const;

  return match(schema)
    .returnType<IRType>()
    .with({ type: "string" }, () => ({
      kind: "string",
      ...common,
    }))
    .with({ type: "boolean" }, () => ({
      kind: "bool",
      ...common,
    }))
    .with({ type: "integer" }, () => ({
      kind: "int",
      ...common,
    }))
    .with({ type: "array", items: P.record(P.string, P.any) }, (schema) => ({
      kind: "array",
      ...common,
      element: toIR({ schema: schema.items, required: null }),
    }))
    .with(
      {
        type: "object",
        properties: P.record(P.string, P.any),
        required: P.array(P.string).optional(),
      },
      (schema) => {
        return {
          kind: "object",
          ...common,
          fields: Object.fromEntries(
            Object.entries(schema.properties).map(
              ([propertyName, propertySchema]) => [
                propertyName,
                toIR({
                  schema: propertySchema,
                  required: schema.required?.includes(propertyName) ?? null,
                }),
              ]
            )
          ),
        };
      }
    )
    .otherwise(() => {
      throw new Error(
        `Unsupported swagger schema type: ${JSON.stringify(schema)}`
      );
    });
}
