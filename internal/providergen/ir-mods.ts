import { match } from "ts-pattern";
import type { IRResource, IRType } from "./ir";
import type { RESOURCES } from "./settings";

/**
 * Converts all arrays of objects from TF attributes to TF blocks.
 *
 * This is a common pattern used in our existing codebase using SDKv2.
 */
function withLegacyBlocks<T extends IRResource | IRType>(ir: T): T {
  return match(ir as IRResource | IRType)
    .with({ kind: "resource" }, { kind: "object" }, (ir) => ({
      ...ir,
      blocks: true,
      fields: Object.fromEntries(
        Object.entries(ir.fields).map(([k, v]) => [k, withLegacyBlocks(v)])
      ),
    }))
    .with({ kind: "array", element: { kind: "object" } }, (ir) => ({
      ...ir,
      element: withLegacyBlocks(ir.element),
      blocks: true,
    }))
    .otherwise((ir) => ir) as T;
}

export const IR_MODS: {
  [K in (typeof RESOURCES)[number]]?: (
    ir: IRResource
  ) => IRResource | Promise<IRResource>;
} = {};
