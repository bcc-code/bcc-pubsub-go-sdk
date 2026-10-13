package model

import (
	"encoding/json"
	"testing"
	"time"
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

func TestPublishRequestDataAttributes(t *testing.T) {
	d := NewPublishRequestData("my.type", "payload")

	if _, ok := marshalToMap(t, d)["attributes"]; ok {
		t.Error("attributes key must be absent when nil")
	}

	d.SetAttributes(map[string]string{"tenant": "abc"})
	attrs, ok := marshalToMap(t, d)["attributes"].(map[string]any)
	if !ok || attrs["tenant"] != "abc" {
		t.Errorf("unexpected attributes: %v", attrs)
	}
}

func TestPublishRequestDataScheduleTime(t *testing.T) {
	d := NewPublishRequestData("my.type", "payload")

	if _, ok := marshalToMap(t, d)["scheduleTime"]; ok {
		t.Error("scheduleTime key must be absent when nil")
	}

	d.SetScheduleTime(time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC))
	if got := marshalToMap(t, d)["scheduleTime"]; got != "2030-01-02T03:04:05Z" {
		t.Errorf("scheduleTime = %v, want 2030-01-02T03:04:05Z", got)
	}
}
