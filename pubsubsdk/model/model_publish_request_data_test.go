package model

import (
	"encoding/json"
	"testing"
)

func marshalToMap(t *testing.T, v any) map[string]any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]any{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestPublishRequestDataMarshalsData(t *testing.T) {
	m := marshalToMap(t, NewPublishRequestData("my.type", "payload"))

	if m["data"] != "payload" {
		t.Errorf("data = %v, want payload", m["data"])
	}
	if _, ok := m["message"]; ok {
		t.Error("message key must not be present")
	}
	if m["type"] != "my.type" {
		t.Errorf("type = %v, want my.type", m["type"])
	}
}
