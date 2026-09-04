import { produce } from "immer";
import type { AttributeListNested, ResourceConfig } from "../schema";
import { findAttribute, findBlock } from "../schema-helpers";

const OPERATOR_ENUM = [
  "is",
  "is_not",
  "contains",
  "does_not_contain",
  "is_one_of",
  "is_not_one_of",
  "is_empty",
  "is_not_empty",
  "contains_key",
  "does_not_contain_key",
  "starts_with",
  "does_not_start_with",
  "matches",
  "does_not_match",
  "is_set",
  "is_not_set",
];

const RULE_TYPE_ENUM = [
  "alert_urgency",
  "working_hour",
  "json_path",
  "field",
  "service",
  "deferral_window",
  "source",
  "related_incidents",
];

const FIELDABLE_TYPE_ENUM = ["AlertField"];

function modifyCondition(condition: AttributeListNested) {
  const ruleType = findAttribute(condition.attributes, "rule_type", "string");
  ruleType.enum = RULE_TYPE_ENUM;

  const fieldableId = findAttribute(
    condition.attributes,
    "fieldable_id",
    "string",
  );
  fieldableId.description =
    "The ID of the alert field. Only used with `field` rule type.";

  const fieldableType = findAttribute(
    condition.attributes,
    "fieldable_type",
    "string",
  );
  fieldableType.description =
    "The type of the fieldable. Only used with `field` rule type.";
  fieldableType.enum = FIELDABLE_TYPE_ENUM;

  const operator = findAttribute(condition.attributes, "operator", "string");
  operator.description =
    "How the value should be matched. For `json_path` rule type: `is`, `is_not`, `contains`, `does_not_contain`. For `field` rule type: `is`, `is_not`, `contains`, `does_not_contain`, `is_one_of`, `is_not_one_of`, `is_empty`, `is_not_empty`, `contains_key`, `does_not_contain_key`, `starts_with`, `does_not_start_with`, `matches`, `does_not_match`. For `source` rule type: `is`, `is_not`, `is_one_of`, `is_not_one_of`. For `related_incidents` rule type: `is_set`, `is_not_set`.";
  operator.enum = OPERATOR_ENUM;

  const serviceIds = findAttribute(condition.attributes, "service_ids", "set");
  serviceIds.description =
    "Service ids for which this escalation path should be used. Only used with `service` rule type.";

  const timeBlocks = findBlock(condition.blocks, "time_blocks", "list_nested");
  timeBlocks.description =
    "Time windows during which alerts are deferred. Only used with `deferral_window` rule type.";

  const values = findAttribute(condition.attributes, "values", "set");
  values.description =
    "Values to match against. Used with `field` and `source` rule types.";
}

export default {
  name: "escalation_path",
  description: "Manages an escalation path.",
  modifyDef: produce((def) => {
    const afterDeferralBehavior = findAttribute(
      def.attributes,
      "after_deferral_behavior",
      "string",
    );
    afterDeferralBehavior.description =
      "What happens after a deferral path finishes. Required for deferral paths.";

    const pathType = findAttribute(def.attributes, "path_type", "string");
    pathType.planModifiers = ["stringplanmodifier.RequiresReplace()"];

    const rules = findBlock(def.blocks, "rules", "list_nested");
    rules.hacks = { nullBehavior: "empty" };

    modifyCondition(rules);

    const timeRestrictions = findBlock(
      def.blocks,
      "time_restrictions",
      "list_nested",
    );
    timeRestrictions.hacks = { nullBehavior: "empty" };

    const notificationTypeRules = findBlock(
      def.blocks,
      "notification_type_rules",
      "list_nested",
    );
    notificationTypeRules.hacks = {
      unknownBehavior: "omit",
      nullBehavior: "omit",
      emptyBehavior: "omit",
    };
    notificationTypeRules.validators = ["listvalidator.SizeAtMost(10)"];

    const notificationTypeRulesConditions = findBlock(
      notificationTypeRules.blocks,
      "conditions",
      "list_nested",
    );
    notificationTypeRulesConditions.validators = [
      "listvalidator.SizeBetween(1, 5)",
    ];

    modifyCondition(notificationTypeRulesConditions);

    const notificationTypeFallback = findAttribute(
      def.attributes,
      "notification_type_fallback",
      "string",
    );
    notificationTypeFallback.default = undefined;
  }),
} satisfies ResourceConfig;
