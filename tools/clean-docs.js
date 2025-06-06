const fs = require("fs");

const file = process.argv[2];
let text = fs.readFileSync(file, "utf-8");

// Clean the Read-Only section
text = text.replace(
  /\nRead-Only\:\n\n- `id` \(String\) The ID of this resource./g,
  "- `id` (String)"
);

// Check if the doc contains 'slug' (case-insensitive)
const useSlug = /slug/i.test(text);
// Replace import block
text = text.replace(
  /import \{\s*\n\s*to ([^\n]+)\.primary\n\s*id = "[^"]+"\s*\n\}/,
  useSlug
    ? 'import {\n  to $1.my-resource\n  slug = "my-resource-slug"\n}'
    : 'import {\n  to $1.my-resource\n  id = "00000000-0000-0000-0000-000000000000"\n}'
);
// Replace terraform import command
text = text.replace(
  /terraform import Resource\.([^\s]+) ([^\n]+)\n/,
  useSlug
    ? 'terraform import $1.my-resource my-resource-slug\n'
    : 'terraform import $1.my-resource 00000000-0000-0000-0000-000000000000\n'
);

fs.writeFileSync(file, text, "utf-8");
