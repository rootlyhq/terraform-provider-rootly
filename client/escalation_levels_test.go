package client

import (
	"io"
	"strings"
	"testing"

	"github.com/rootlyhq/terraform-provider-rootly/v5/tools"
)

// Delay is a nullable *int precisely so these three cases stay distinguishable on the wire:
//
//   - nil   -> attribute omitted, server keeps whatever delay it already has. Needed so an
//     update that only touches notification_target_params can't reset delay to 0.
//   - 0     -> attribute sent as 0, server applies it. Needed so `delay = 0` is settable
//     without a permadiff (TER-182, #351).
//   - n     -> attribute sent as n.
//
// A plain int can satisfy at most two of the three, which is why this file exists.
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

// Regression guard for TER-182 (#351): an explicitly requested delay of 0 must still reach
// the API, otherwise `delay = 0` silently no-ops and the plan never converges.
func TestEscalationLevelDelaySentWhenExplicitlyZero(t *testing.T) {
	level := &EscalationLevel{
		ID:    "123",
		Delay: tools.Int(0),
	}

	body := marshalEscalationLevel(t, level)
	if !strings.Contains(body, `"delay":0`) {
		t.Fatalf("expected an explicit delay of 0 to be present in payload, got: %s", body)
	}
}

func TestEscalationLevelDelaySentWhenSet(t *testing.T) {
	level := &EscalationLevel{
		ID:    "123",
		Delay: tools.Int(15),
	}

	body := marshalEscalationLevel(t, level)
	if !strings.Contains(body, `"delay":15`) {
		t.Fatalf("expected delay to be present and set to 15 in payload, got: %s", body)
	}
}

// The read path has to survive the pointer too: a delay in the response must come back as a
// non-nil pointer, and an absent one must stay nil rather than erroring.
func TestEscalationLevelDelayUnmarshal(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		lvl := unmarshalEscalationLevel(t, `{"data":{"type":"escalation_levels","id":"123","attributes":{"delay":15,"position":1}}}`)
		if lvl.Delay == nil {
			t.Fatal("expected Delay to be non-nil")
		}
		if *lvl.Delay != 15 {
			t.Fatalf("expected Delay 15, got %d", *lvl.Delay)
		}
	})

	t.Run("absent", func(t *testing.T) {
		lvl := unmarshalEscalationLevel(t, `{"data":{"type":"escalation_levels","id":"123","attributes":{"position":1}}}`)
		if lvl.Delay != nil {
			t.Fatalf("expected Delay to be nil, got %d", *lvl.Delay)
		}
	})
}

func unmarshalEscalationLevel(t *testing.T, payload string) *EscalationLevel {
	t.Helper()

	out, err := UnmarshalData(io.NopCloser(strings.NewReader(payload)), new(EscalationLevel))
	if err != nil {
		t.Fatalf("UnmarshalData returned error: %s", err)
	}

	lvl, ok := out.(*EscalationLevel)
	if !ok {
		t.Fatalf("UnmarshalData returned %T, want *EscalationLevel", out)
	}

	return lvl
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
