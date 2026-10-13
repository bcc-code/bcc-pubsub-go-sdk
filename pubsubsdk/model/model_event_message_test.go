package model

import (
	"encoding/json"
	"testing"
)

func TestEventMessageUnmarshalAttributes(t *testing.T) {
	body := `{"id":"1","source":"s","specversion":"1.0","type":"t","data":"{}","time":"2026-01-01T00:00:00Z","attributes":{"tenant":"abc"}}`

	var m EventMessage
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Fatal(err)
	}

	if m.Attributes["tenant"] != "abc" {
		t.Errorf("Attributes = %v, want tenant=abc", m.Attributes)
	}
	if m.GetId() != "1" || m.GetSpecVersion() != "1.0" {
		t.Errorf("unexpected message: id=%q specVersion=%q", m.GetId(), m.GetSpecVersion())
	}
}

func TestEventMessageMarshalAttributes(t *testing.T) {
	m := NewEventMessage()

	if _, ok := marshalToMap(t, m)["attributes"]; ok {
		t.Error("attributes key must be absent when nil")
	}

	m.SetAttributes(map[string]string{"tenant": "abc"})
	attrs, ok := marshalToMap(t, m)["attributes"].(map[string]any)
	if !ok || attrs["tenant"] != "abc" {
		t.Errorf("unexpected attributes: %v", attrs)
	}
}
