package diffsuppressfunc

import "testing"

func TestNormalizeTimeString(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		// 24-hour format
		{"08:00", "08:00"},
		{"17:30", "17:30"},
		{"23:59", "23:59"},
		// 12-hour with AM/PM, various spacings
		{"08:00 AM", "08:00"},
		{"8:00 AM", "08:00"},
		{"08:00AM", "08:00"},
		{"8:00AM", "08:00"},
		{"5:15 PM", "17:15"},
		{"05:15 PM", "17:15"},
		{"5:15PM", "17:15"},
		{"05:15PM", "17:15"},
		{"12:00PM", "12:00"},
		{"12:00 AM", "00:00"},
		// Lower/upper case tolerance
		{"8:00 am", "08:00"},
		{"5:15 pm", "17:15"},
		// Whitespace
		{"   09:00 AM  ", "09:00"},
		{"  14:00 ", "14:00"},
		// Fallback (uncleanable string, returns as is)
		{"foo", "FOO"},
		{"", ""},
	}

	for _, c := range cases {
		got := normalizeTimeString(c.input)
		if got != c.expected {
			t.Errorf("normalizeTimeString(%q) == %q, want %q", c.input, got, c.expected)
		}
	}
}
