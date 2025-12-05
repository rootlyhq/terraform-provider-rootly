import type { IRResource } from "./ir";
import type { RESOURCES } from "./settings";

export const IR_MODS: {
  [K in (typeof RESOURCES)[number]]?: (
    ir: IRResource
  ) => IRResource | Promise<IRResource>;
} = {
  dashboard_panel: (ir) => {
    return ir;
  },
};
