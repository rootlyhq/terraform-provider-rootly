package client

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

func marshalEscalationPathAttrs(t *testing.T, p *EscalationPath) map[string]interface{} {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, p); err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Data struct {
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	return payload.Data.Attributes
}

func TestEscalationPathRetriggerTimeoutMarshal(t *testing.T) {
	// nil pointer (cleared / not configured): the attribute must be present as
	// JSON null so the API resets it to inherit, NOT omitted.
	attrs := marshalEscalationPathAttrs(t, &EscalationPath{ID: "ep1", Name: "renamed"})
	if v, present := attrs["retrigger_timeout_minutes"]; !present || v != nil {
		t.Fatalf("expected retrigger_timeout_minutes present and null, got %v (present=%v)", v, present)
	}

	// -1 (never): value must be serialized.
	never := -1
	attrs = marshalEscalationPathAttrs(t, &EscalationPath{ID: "ep1", RetriggerTimeoutMinutes: &never})
	if v := attrs["retrigger_timeout_minutes"]; v != float64(-1) {
		t.Fatalf("expected -1, got %v", v)
	}

	// positive minutes: value must be serialized.
	thirty := 30
	attrs = marshalEscalationPathAttrs(t, &EscalationPath{ID: "ep1", RetriggerTimeoutMinutes: &thirty})
	if v := attrs["retrigger_timeout_minutes"]; v != float64(30) {
		t.Fatalf("expected 30, got %v", v)
	}
}

func TestEscalationPathRetriggerTimeoutUnmarshalNull(t *testing.T) {
	body := `{"data":{"id":"ep1","type":"escalation_paths","attributes":{"name":"n","retrigger_timeout_minutes":null}}}`
	p := new(EscalationPath)
	if err := jsonapi.UnmarshalPayload(strings.NewReader(body), p); err != nil {
		t.Fatal(err)
	}
	if p.RetriggerTimeoutMinutes != nil {
		t.Fatalf("expected nil retrigger timeout, got %d", *p.RetriggerTimeoutMinutes)
	}
}
