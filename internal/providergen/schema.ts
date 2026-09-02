import type { oas30 } from "openapi3-ts";

export interface ClientConfig {
  name: string;
  actions?: {
    list?: {
      enabled: true;
    };
    get?: {
      enabled: true;
    };
    create?: {
      enabled: true;
    };
    update?: {
      enabled: true;
    };
  };
}

interface DataSourceListConfig {
  strategy: "list";
  resourceName: string;
}

interface DataSourceSingleConfig {
  strategy: "single";
}

export type DataSourceConfig = {
  name: string;
  description?: string;
} & (DataSourceListConfig | DataSourceSingleConfig);

export type ResolvedDataSourceConfig = {
  type: "data_source";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the data source. */
    struct: `${string}DataSource`;
    /** Name of the struct that represents the model of the data source. */
    model: `${string}DataSourceModel`;
  };
  schemas: {
    read: oas30.SchemaObject;
    resolved: oas30.SchemaObject;
  };
} & (DataSourceListConfig | DataSourceSingleConfig);

// TODO
export type ResourceConfig = {
  name: string;
  description?: string;
};

export type ResolvedResourceConfig = {
  type: "resource";
  name: string;
  description?: string;
  goNames: {
    /** Name of the struct that represents the resource. */
    struct: `${string}Resource`;
    /** Name of the struct that represents the model of the resource. */
    model: `${string}ResourceModel`;
  };
  schemas: {
    create: oas30.SchemaObject;
    read: oas30.SchemaObject;
    update: oas30.SchemaObject;
    resolved: oas30.SchemaObject;
  };
};

export type ComputedOptionalRequired =
  | "computed"
  | "optional"
  | "computed_optional"
  | "required";

declare module "openapi3-ts/oas30" {
  interface SchemaObject {
    /** Defines the jsonapi attribute tag. */
    "x-go-jsonapi-tag"?: "primary" | "attr";
    /** Defines the jsonapi type. */
    "x-go-jsonapi-type"?: string;
    /** Overrides the nested object type name. Only used for nested objects. */
    "x-go-nested-type"?: string;
    /** Defines if the property is computed, optional or required. */
    "x-tf-computed-optional-required"?: ComputedOptionalRequired;
    /** Indicates this is the top level item type for the data source. Only used for plural data sources. */
    "x-tf-top-level-item-type"?: boolean;
    /** The type of collection to use for arrays. */
    "x-tf-collection-type"?: "list" | "set";
  }
}
