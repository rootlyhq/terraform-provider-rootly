package client

import (
	"io"
	"strings"
	"testing"
)

// Regression test for a bug where partial updates (e.g. only changing
// notification_target_params) serialized delay:0 into the request body even
// though delay was never set by the caller, silently resetting it server-side.
// Every other optional attribute on EscalationLevel already uses omitempty;
// delay must behave the same way.
func TestEscalationLevelDelayOmittedWhenUnset(t *testing.T) {
	level := &EscalationLevel{
		ID: "123",
		NotificationTargetParams: []interface{}{
			map[string]interface{}{"type": "schedule", "id": "sched-1"},
		},
	}

	body := marshalEscalationLevel(t, level)
	if strings.Contains(body, `"delay"`) {
		t.Fatalf("expected delay to be omitted from payload when unset, got: %s", body)
	}
}

func TestEscalationLevelDelaySentWhenSet(t *testing.T) {
	level := &EscalationLevel{
		ID:    "123",
		Delay: 15,
	}

	body := marshalEscalationLevel(t, level)
	if !strings.Contains(body, `"delay":15`) {
		t.Fatalf("expected delay to be present and set to 15 in payload, got: %s", body)
	}
}

func marshalEscalationLevel(t *testing.T, level *EscalationLevel) string {
	t.Helper()

	reader, err := MarshalData(level)
	if err != nil {
		t.Fatalf("MarshalData returned error: %s", err)
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed reading marshaled payload: %s", err)
	}

	return string(b)
}
