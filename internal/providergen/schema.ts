import z from "zod";

const ParameterObject = z.looseObject({});

const OperationObject = z.looseObject({
  parameters: z.array(ParameterObject).optional(),
  operationId: z.string(),
});

const PathItemObject = z.strictObject({
  parameters: z.array(ParameterObject).optional(),
  // HTTP methods
  get: OperationObject.optional(),
  put: OperationObject.optional(),
  post: OperationObject.optional(),
  delete: OperationObject.optional(),
  patch: OperationObject.optional(),
});

const SchemaObject = z.lazy(() =>
  z.object({
    type: z.enum(["object", "array", "string", "integer", "boolean"]),
    description: z.string().optional(),
    // properties: z.record(z.string(), SchemaObject).optional(),
    required: z.array(z.string()).optional(),
  })
);

const TopLevelSchemaObject = z.strictObject({
  type: z.literal("object"),
  description: z.string().optional(),
  properties: z.record(z.string(), SchemaObject),
  required: z.array(z.string()).optional(),
});

export const OpenAPIObject = z.object({
  paths: z.record(z.string(), PathItemObject),
  components: z.object({
    schemas: z.record(z.string(), TopLevelSchemaObject),
  }),
});
export type OpenAPIObject = z.infer<typeof OpenAPIObject>;
