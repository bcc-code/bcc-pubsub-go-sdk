package model

import "testing"

func TestSubscribeRequestDataOrgs(t *testing.T) {
	d := NewSubscribeRequestData("my.type", "https://example.com/hook")

	if _, ok := marshalToMap(t, d)["orgs"]; ok {
		t.Error("orgs key must be absent when nil")
	}

	d.SetOrgs([]string{"org-a", "org-b"})
	orgs, ok := marshalToMap(t, d)["orgs"].([]any)
	if !ok || len(orgs) != 2 || orgs[0] != "org-a" || orgs[1] != "org-b" {
		t.Errorf("unexpected orgs: %v", orgs)
	}
}
