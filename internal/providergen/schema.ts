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

export const OpenAPIObject = z.object({
  paths: z.record(z.string(), PathItemObject),
  components: z.object({
    schemas: z.record(z.string(), z.any()),
  }),
});
export type OpenAPIObject = z.infer<typeof OpenAPIObject>;
