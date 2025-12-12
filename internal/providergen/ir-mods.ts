import type { ResourceIR } from "./ir";

export const IR_MODS: Record<
  string,
  (ir: ResourceIR) => ResourceIR | Promise<ResourceIR>
> = {};
