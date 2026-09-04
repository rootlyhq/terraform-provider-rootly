import { produce } from "immer";
import type { ResourceConfig } from "../schema";
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

export default {
  name: "escalation_path",
  description: "Manages an escalation path.",
  modifyDef: produce((def) => {
    const pathType = findAttribute(def.attributes, "path_type", "string");
    pathType.planModifiers = ["stringplanmodifier.RequiresReplace()"];

    const rules = findBlock(def.blocks, "rules", "list_nested");
    rules.hacks = { nullBehavior: "empty" };

    const rulesOperator = findAttribute(rules.attributes, "operator", "string");
    rulesOperator.enum = OPERATOR_ENUM;

    const rulesRuleType = findAttribute(
      rules.attributes,
      "rule_type",
      "string",
    );
    rulesRuleType.enum = RULE_TYPE_ENUM;

    const rulesFieldableType = findAttribute(
      rules.attributes,
      "fieldable_type",
      "string",
    );
    rulesFieldableType.enum = FIELDABLE_TYPE_ENUM;

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

    const notificationTypeRulesConditionsOperator = findAttribute(
      notificationTypeRulesConditions.attributes,
      "operator",
      "string",
    );
    notificationTypeRulesConditionsOperator.enum = OPERATOR_ENUM;

    const notificationTypeRulesConditionsRuleType = findAttribute(
      notificationTypeRulesConditions.attributes,
      "rule_type",
      "string",
    );
    notificationTypeRulesConditionsRuleType.enum = RULE_TYPE_ENUM;

    const notificationTypeRulesConditionsFieldableType = findAttribute(
      notificationTypeRulesConditions.attributes,
      "fieldable_type",
      "string",
    );
    notificationTypeRulesConditionsFieldableType.enum = FIELDABLE_TYPE_ENUM;

    const notificationTypeFallback = findAttribute(
      def.attributes,
      "notification_type_fallback",
      "string",
    );
    notificationTypeFallback.default = undefined;
  }),
} satisfies ResourceConfig;
