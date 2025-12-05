import { generateResourceIR } from "./ir";
import { SWAGGER_MODS } from "./swagger-mods";
import { IR_MODS } from "./ir-mods";
import { RESOURCES } from "./settings";
import {
  generateClientModelFile,
  generateProvider,
  generateResource,
} from "./generators";

async function getRootlySwagger() {
  const response = await fetch(
    "https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json"
  );
  return (await response.json()) as any;
}

async function writeAndFormatGoFile(destination: URL, code: string) {
  await Bun.write(destination, code);
  await Bun.$`go fmt ${destination.pathname}`;
  await Bun.$`go tool goimports -w ${destination.pathname}`;
}

async function main() {
  console.log("🚀 Fetching Rootly Swagger...");
  let swagger = await getRootlySwagger();
  await Bun.write(
    new URL("out/swagger.original.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  console.log("🚀 Modifying Rootly Swagger...");
  for (const mod of SWAGGER_MODS) {
    swagger = await mod(swagger);
  }
  await Bun.write(
    new URL("out/swagger.modified.json", import.meta.url),
    JSON.stringify(swagger, null, 2)
  );

  for (const name of RESOURCES) {
    console.log(`🚀 Generating ${name}...`);

    let ir = generateResourceIR({ swagger, name });
    await Bun.write(
      new URL(`out/ir_${name}.original.json`, import.meta.url),
      JSON.stringify(ir, null, 2)
    );

    if (IR_MODS[name]) {
      ir = await IR_MODS[name](ir);
    }
    await Bun.write(
      new URL(`out/ir_${name}.modified.json`, import.meta.url),
      JSON.stringify(ir, null, 2)
    );

    // Client
    {
      const code = generateClientModelFile({ ir, name });

      await writeAndFormatGoFile(
        new URL(`../apiclient/model_resource_${name}_gen.go`, import.meta.url),
        code
      );
    }

    // Resource
    {
      const code = generateResource({ ir, name });

      await writeAndFormatGoFile(
        new URL(`../provider/resource_${name}_gen.go`, import.meta.url),
        code
      );
    }
  }

  // Provider
  {
    const code = generateProvider({ resources: RESOURCES });

    await writeAndFormatGoFile(
      new URL(`../provider/provider_gen.go`, import.meta.url),
      code
    );
  }
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
