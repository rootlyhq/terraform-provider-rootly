import { $ } from "bun";
import type {
  MySchema,
  Provider,
  Resource,
  ResourceAttributes,
  ResourceBoolAttribute,
  ResourceInt64Attribute,
  ResourceSingleNestedAttribute,
  ResourceStringAttribute,
  SchemaCustomType,
} from "./tfplugingen-framework-types";
import $RefParser from "@apidevtools/json-schema-ref-parser";
import { type OpenAPIObject, type SchemaObject } from "openapi3-ts/oas30";
import {
  renameTfKeysDeep,
  unwrapJsonApi,
  unwrapSingleAllOfDeep,
} from "./utils";

const specPath = new URL("./specification.json", import.meta.url);
const providerPath = new URL("../provider", import.meta.url);

const resources = ["override_shift"];
const excludedProperties = ["created_at", "updated_at"];

interface Hints {
  requiredFields: string[];
}

async function getRootlySwagger() {
  const response = await fetch(
    "https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json"
  );
  let data = (await response.json()) as OpenAPIObject;
  data = renameTfKeysDeep(data);
  data = unwrapSingleAllOfDeep(data);
  data = await $RefParser.dereference(data);
  return data;
}

function buildTerraformSpec({ swagger }: { swagger: OpenAPIObject }) {
  const providerSpec: Provider = {
    name: "rootly",
    schema: {
      attributes: [
        {
          name: "api_host",
          string: {
            description:
              "The Rootly API host. Defaults to https://api.rootly.com. Can also be sourced from the `ROOTLY_API_URL` environment variable.",
            optional_required: "optional",
          },
        },
        {
          name: "api_token",
          string: {
            description:
              "The Rootly API Token. Generate it from your account at https://rootly.com/account. It must be provided but can also be sourced from the `ROOTLY_API_TOKEN` environment variable.",
            optional_required: "optional",
            sensitive: true,
          },
        },
      ],
    },
  };

  const spec: MySchema = {
    $schema:
      "https://raw.githubusercontent.com/hashicorp/terraform-plugin-codegen-spec/main/spec/v0.1/schema.json",
    version: "0.1",
    provider: providerSpec,
    resources: resources.map((name) =>
      buildTerraformResourceSpec({ swagger, name })
    ),
    data_sources: [],
  };

  return spec;
}

function buildTerraformResourceSpec({
  swagger,
  name,
}: {
  swagger: any;
  name: string;
}): Resource {
  const component = swagger.components.schemas[name] as {
    properties: Record<string, any>;
  };

  const requiredFields =
    swagger.components.schemas[`new_${name}`].properties.data.properties
      .attributes.required;
  const hints: Hints = {
    requiredFields,
  };

  return {
    name,
    schema: {
      attributes: [
        {
          name: "id",
          string: {
            description: "The unique identifier of this resource.",
            computed_optional_required: "computed",
          },
        },
        ...Object.entries(component.properties)
          .map(([name, propertySchema]) =>
            buildTerraformResourceAttributeSpec({
              swagger,
              name,
              propertySchema,
              hints,
            })
          )
          .filter((attr) => attr !== null),
      ],
    },
  };
}

function buildTerraformResourceAttributeSpec({
  swagger,
  name,
  propertySchema,
  hints,
}: {
  swagger: OpenAPIObject;
  name: string;
  propertySchema: any;
  hints?: Hints;
}): ResourceAttributes[number] | null {
  if (excludedProperties.includes(name) || propertySchema["x-rootly-ignore"]) {
    return null;
  }

  const computed_optional_required = hints?.requiredFields.includes(name)
    ? "required"
    : "optional";

  switch (propertySchema.type) {
    case "string":
      let custom_type: SchemaCustomType | undefined;
      if (name.endsWith("_at")) {
        custom_type = {
          import: {
            path: "github.com/rootlyhq/terraform-provider-rootly/v2/internal/rootlytypes",
          },
          type: "rootlytypes.RFC3339Type{}",
          value_type: "rootlytypes.RFC3339",
        };
      }

      return {
        name,
        string: {
          description: propertySchema.description,
          computed_optional_required,
          custom_type,
        },
      } as ResourceStringAttribute;

    case "boolean":
      return {
        name,
        bool: {
          description: propertySchema.description,
          computed_optional_required,
        },
      } as ResourceBoolAttribute;

    case "integer":
      return {
        name,
        int64: {
          description: propertySchema.description,
          computed_optional_required,
        },
      } as ResourceInt64Attribute;

    case "object":
      if (
        "data" in propertySchema.properties &&
        "id" in propertySchema.properties.data.properties &&
        "type" in propertySchema.properties.data.properties &&
        "attributes" in propertySchema.properties.data.properties
      ) {
        const attributes = Object.entries(
          propertySchema.properties.data.properties.attributes.properties
        )
          .map(([name, attributeSchema]: [string, any]) =>
            buildTerraformResourceAttributeSpec({
              swagger,
              name,
              propertySchema: attributeSchema,
            })
          )
          .filter((attr) => attr !== null);

        return {
          name,
          single_nested: {
            description: propertySchema.description,
            computed_optional_required,
            attributes: [
              ...("id" in
              propertySchema.properties.data.properties.attributes.properties
                ? []
                : [
                    buildTerraformResourceAttributeSpec({
                      swagger,
                      name: "id",
                      propertySchema:
                        propertySchema.properties.data.properties.id,
                    }),
                  ]),
              ...attributes,
            ].filter((attr) => attr !== null),
          },
        } as ResourceSingleNestedAttribute;
      } else {
        const attributes = Object.entries(propertySchema.properties)
          .map(([name, attributeSchema]: [string, any]) =>
            buildTerraformResourceAttributeSpec({
              swagger,
              name,
              propertySchema: attributeSchema,
            })
          )
          .filter((attr) => attr !== null);

        return {
          name,
          single_nested: {
            description: propertySchema.description,
            computed_optional_required,
            attributes,
          },
        } as ResourceSingleNestedAttribute;
      }

    default:
      throw new Error(`Unsupported property type: ${propertySchema.type}`);
  }
}

async function main() {
  console.log("🚀 Fetching Rootly Swagger...");
  const swagger = await getRootlySwagger();

  await Bun.write(
    new URL("swagger.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  console.log("🚀 Generating Terraform spec...");
  const terraformSpec = buildTerraformSpec({ swagger });

  console.log("🚀 Writing Terraform spec...");
  await Bun.write(specPath, JSON.stringify(terraformSpec, null, 2));

  console.log("🚀 Generating Terraform provider...");
  await $`go tool tfplugingen-framework generate all --input ${specPath.pathname} --output ${providerPath.pathname}`;
}

await main()
  .then(() => {
    console.log("✨ Done");
    process.exit(0);
  })
  .catch((e) => {
    console.error(e);
    process.exit(1);
  });
