const inflect = require("inflect");

module.exports = {
  ...inflect,
  camelize: (name) => {
    if (name.match(/^post_mortem_templates?$/)) {
      return name.replace("post_mortem_template", "PostmortemTemplate");
    }
    if (name.match(/^slas?$/)) {
      return name.replace(/^sla/, "SLA");
    }
    return inflect.camelize(name);
  },
};
