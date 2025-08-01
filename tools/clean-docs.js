const fs = require("fs");

const file = process.argv[2];
let text = fs.readFileSync(file, "utf-8");

// Clean the Read-Only section
text = text.replace(
  /\nRead-Only\:\n\n- `id` \(String\) The ID of this resource./g,
  "- `id` (String)"
);

// Replace import block
text = text.replace(
  /import \{\s*\n\s*to ([^\n]+)\.primary\n\s*id = "[^"]+"\s*\n\}/,
    'import {\n  to $1.primary\n  id = "a816421c-6ceb-481a-87c4-585e47451f24"\n}'
);
// Replace terraform import command
text = text.replace(
  /terraform import Resource\.([^\s]+) ([^\n]+)\n/,
    'terraform plan -generate-config-out=generated.tf'
);

fs.writeFileSync(file, text, "utf-8");
