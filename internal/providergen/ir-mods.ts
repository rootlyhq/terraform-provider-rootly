import type { IRArray, IRResource } from "./ir";
import type { RESOURCES } from "./settings";

export const IR_MODS: {
  [K in (typeof RESOURCES)[number]]?: (
    ir: IRResource
  ) => IRResource | Promise<IRResource>;
} = {
  alert_route: (ir) => {
    (ir.fields["alerts_source_ids"] as IRArray).distinct = true;
    (ir.fields["owning_team_ids"] as IRArray).distinct = true;
    (ir.fields["rules"] as IRArray).distinct = true;
    return ir;
  },
};
