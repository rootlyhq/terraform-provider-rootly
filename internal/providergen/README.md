# Next-gen provider generator

To install dependencies:

```bash
# At ./internal/providergen
bun install
```

You will need to pull and clean the swagger file first by running the existing
codegen:

```bash
# At root of the provider
make codegen
```

Then generate the provider code using the new generator:

```bash
# At ./internal/providergen
bun run index.ts --input ../../schema/swagger.json
```
