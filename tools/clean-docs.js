const fs = require("fs");

const file = process.argv[2];
const text = fs.readFileSync(file, "utf-8");
fs.writeFileSync(
  file,
  text.replace(
    /\nRead-Only\:\n\n- `id` \(String\) The ID of this resource./g,
    "- `id` (String)"
  ),
  "utf-8"
);
