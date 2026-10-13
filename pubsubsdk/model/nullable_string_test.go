package model

import "testing"

func TestNullableStringSetNil(t *testing.T) {
	s := NullableString("value")
	s.Set(nil)
	if s != "" {
		t.Fatalf("expected value to be unset, got %q", string(s))
	}

	v := "x"
	s.Set(&v)
	if s != "x" {
		t.Fatalf("expected x, got %q", string(s))
	}
}
