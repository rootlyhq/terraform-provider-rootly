import type { AttributeBlockType, AttributeType } from "./schema";

export function findBlock<T extends AttributeBlockType["type"]>(
  blocks: AttributeBlockType[],
  name: string,
  type: T,
): Extract<AttributeBlockType, { type: T }> {
  const block = blocks.find((v) => v.name === name);
  if (!block) {
    throw new Error(`Block "${name}" not found`);
  }
  if (block.type !== type) {
    throw new Error(
      `Block "${name}" has type "${block.type}", expected "${type}"`,
    );
  }
  return block as Extract<AttributeBlockType, { type: T }>;
}

export function findAttribute<T extends AttributeType["type"]>(
  attributes: AttributeType[],
  name: string,
  type: T,
): Extract<AttributeType, { type: T }> {
  const attr = attributes.find((v) => v.name === name);
  if (!attr) {
    throw new Error(`Attribute "${name}" not found`);
  }
  if (attr.type !== type) {
    throw new Error(
      `Attribute "${name}" has type "${attr.type}", expected "${type}"`,
    );
  }
  return attr as Extract<AttributeType, { type: T }>;
}
