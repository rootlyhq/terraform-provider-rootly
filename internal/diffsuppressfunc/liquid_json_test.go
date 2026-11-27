package diffsuppressfunc

import "testing"

func TestLiquidJSON(t *testing.T) {
	cases := []struct {
		name, old, new string
		want           bool
	}{
		{"equal JSON, different formatting", "[\n  {\n    \"type\": \"divider\"\n  }\n]", `[{"type":"divider"}]`, true},
		{"different JSON", `[{"type":"divider"}]`, `[]`, false},
		{"liquid, trailing newline from heredoc", `[{"a":1}{% if x %},{"b":2}{% endif %}]`, "[{\"a\":1}{% if x %},{\"b\":2}{% endif %}]\n", true},
		{"liquid, leading whitespace", `{{ x }}`, "  {{ x }}", true},
		{"liquid, different content", `{{ x }}`, `{{ y }}`, false},
		{"liquid, inner whitespace differs", `{% if x %} []{% endif %}`, `{% if x %}[]{% endif %}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := LiquidJSON("", tc.old, tc.new, nil); got != tc.want {
				t.Errorf("LiquidJSON(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.want)
			}
		})
	}
}
