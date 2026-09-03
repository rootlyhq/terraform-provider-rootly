import { dereferenceSync } from "@trojs/openapi-dereference";
import { oas30 } from "openapi3-ts";
import { parseArgs } from "util";
import { generateClient } from "./generate-client";
import { generateDataSource } from "./generate-data-source";
import { generateProvider } from "./generate-provider";
import { CLIENTS, DATA_SOURCES, RESOURCES } from "./settings";
import { generateDataSourceDef, generateResourceDef } from "./schema";
import { generateResource } from "./generate-resource";

async function parseArguments() {
  const { values } = parseArgs({
    args: Bun.argv,
    options: {
      input: {
        type: "string",
        short: "i",
      },
    },
    strict: true,
    allowPositionals: true,
  });

  if (!values.input) {
    throw new Error("Missing required argument: --input (-i)");
  }

  const file = Bun.file(values.input);
  const exists = await file.exists();
  if (!exists) {
    throw new Error(`File not found: "${values.input}"`);
  }

  let doc: oas30.OpenAPIObject;
  try {
    doc = await file.json();
  } catch {
    throw new Error(`Invalid JSON format in file: "${values.input}"`);
  }

  doc = dereferenceSync(doc);

  return {
    doc,
  };
}

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  const { doc } = await parseArguments();

  for (const config of CLIENTS) {
    const code = generateClient({ doc, config });
    await writeAndFormatGoFile(
      new URL(`../apiclient/${config.name}_gen.go`, import.meta.url),
      code,
    );
  }

  for (const config of DATA_SOURCES) {
    const def = generateDataSourceDef({ doc, config });
    await Bun.write(
      new URL(`dist/data_source_def_${config.name}.json`, import.meta.url),
      JSON.stringify(def, null, 2),
    );

    // const resolvedConfig = resolveDataSourceConfig({ doc, config });
    // const code = generateDataSource({ doc, config: resolvedConfig });
    // await writeAndFormatGoFile(
    //   new URL(
    //     `../provider/data_source_${resolvedConfig.name}.gen.go`,
    //     import.meta.url,
    //   ),
    //   code,
    // );
  }

  for (const config of RESOURCES) {
    const def = generateResourceDef({ doc, config });
    await Bun.write(
      new URL(`dist/resource_def_${config.name}.json`, import.meta.url),
      JSON.stringify(def, null, 2),
    );

    const code = generateResource({ doc, def });
    await writeAndFormatGoFile(
      new URL(`../provider/resource_${def.name}.gen.go`, import.meta.url),
      code,
    );
  }

  {
    const code = generateProvider({
      dataSources: DATA_SOURCES,
      resources: RESOURCES,
    });
    await writeAndFormatGoFile(
      new URL(`../provider/provider.gen.go`, import.meta.url),
      code,
    );
  }
}

await main()
  .then(() => {
    console.log("✨ Done");
    process.exit(0);
  })
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
