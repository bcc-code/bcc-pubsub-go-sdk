package pubsubsdk

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return New(server.URL)
}

func TestRequestProblemDetails(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"https://example.com/bad","title":"Bad request","status":400,"detail":"orgs not in token"}`))
	})

	err := client.Request(context.Background(), http.MethodPost, client.URL("Publish"), nil, nil)

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 400 || apiErr.Title != "Bad request" || apiErr.Detail != "orgs not in token" || apiErr.Type != "https://example.com/bad" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
	if want := "pubsub: 400 Bad Request: orgs not in token"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestRequestValidationProblemDetails(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"https://tools.ietf.org/html/rfc9110#section-15.5.1","title":"One or more validation errors occurred.","status":400,"errors":{"Type":["The Type field is required."],"$.time":["The JSON value could not be converted.","Second message."]}}`))
	})

	err := client.Request(context.Background(), http.MethodPost, client.URL("Publish"), nil, nil)

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if len(apiErr.Errors) != 2 || apiErr.Errors["Type"][0] != "The Type field is required." {
		t.Errorf("unexpected Errors: %+v", apiErr.Errors)
	}
	want := "pubsub: 400 Bad Request: $.time: The JSON value could not be converted., Second message.; Type: The Type field is required."
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestRequestPlainTextError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("no such subscription"))
	})

	err := client.Request(context.Background(), http.MethodGet, client.URL("x"), nil, nil)

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 404 || apiErr.Body != "no such subscription" || apiErr.Detail != "" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
	if want := "pubsub: 404 Not Found: no such subscription"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestRequestEmptyBodyError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	err := client.Request(context.Background(), http.MethodGet, client.URL("x"), nil, nil)

	if want := "pubsub: 401 Unauthorized"; err == nil || err.Error() != want {
		t.Errorf("Error() = %v, want %q", err, want)
	}
}

func TestRequestTruncatesBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(make([]byte, 5000))
	})

	err := client.Request(context.Background(), http.MethodGet, client.URL("x"), nil, nil)

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if len(apiErr.Body) != 1024 {
		t.Errorf("len(Body) = %d, want 1024", len(apiErr.Body))
	}
}

func TestRequestDecodesResult(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"value"}`))
	})

	var result struct {
		Name string `json:"name"`
	}
	if err := client.Request(context.Background(), http.MethodGet, client.URL("x"), nil, &result); err != nil {
		t.Fatal(err)
	}
	if result.Name != "value" {
		t.Errorf("Name = %q, want value", result.Name)
	}
}
