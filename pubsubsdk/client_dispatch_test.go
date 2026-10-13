package pubsubsdk

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

func TestDispatchPostsPayloadAndDecodesResponse(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"dispatchId":"d1","endpoint":"https://example.com/hook","scheduleTime":"2030-01-02T03:04:05+01:00"}`))
	}))
	defer server.Close()

	scheduleTime := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	response, err := New(server.URL).Dispatch(context.Background(), model.DispatchPayload{
		DispatchId:       "d1",
		Type:             "my.type",
		Method:           "POST",
		Endpoint:         "https://example.com/hook",
		ScheduleTime:     scheduleTime,
		DeliverUnwrapped: true,
		Data:             map[string]any{"a": 1},
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/Dispatch" {
		t.Errorf("path = %s, want /Dispatch", gotPath)
	}
	want := map[string]any{
		"dispatchId":       "d1",
		"type":             "my.type",
		"method":           "POST",
		"endpoint":         "https://example.com/hook",
		"scheduleTime":     "2030-01-02T03:04:05Z",
		"deliverUnwrapped": true,
		"data":             map[string]any{"a": float64(1)},
	}
	gotJSON, _ := json.Marshal(gotBody)
	wantJSON, _ := json.Marshal(want)
	if string(gotJSON) != string(wantJSON) {
		t.Errorf("body = %s, want %s", gotJSON, wantJSON)
	}

	if response.DispatchId != "d1" || response.Endpoint != "https://example.com/hook" {
		t.Errorf("unexpected response: %+v", response)
	}
	if wantTime := time.Date(2030, 1, 2, 2, 4, 5, 0, time.UTC); !response.ScheduleTime.Equal(wantTime) {
		t.Errorf("scheduleTime = %v, want %v", response.ScheduleTime, wantTime)
	}
}

func TestDispatchOmitsNilData(t *testing.T) {
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	_, err := New(server.URL).Dispatch(context.Background(), model.DispatchPayload{DispatchId: "d1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := gotBody["data"]; ok {
		t.Error("data key must be absent when nil")
	}
	if _, ok := gotBody["scheduleTime"]; !ok {
		t.Error("scheduleTime is required and must be present")
	}
}

func TestDispatchReturnsProblemDetailsError(t *testing.T) {
	server := newProblemServer(t)
	defer server.Close()

	response, err := New(server.URL).Dispatch(context.Background(), model.DispatchPayload{DispatchId: "d1"})
	if response != nil {
		t.Errorf("response = %+v, want nil", response)
	}
	assertBadRequestError(t, err)
}

func TestDeleteTask(t *testing.T) {
	var gotMethod, gotURI string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotURI = r.RequestURI
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Task deleted successfully"))
	}))
	defer server.Close()

	if err := New(server.URL).DeleteTask(context.Background(), "id 1&x"); err != nil {
		t.Fatal(err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if want := "/Dispatch/Delete?dispatchId=id+1%26x"; gotURI != want {
		t.Errorf("request URI = %s, want %s", gotURI, want)
	}
}

func TestDeleteTaskEmptyIDSendsNothing(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	if err := New(server.URL).DeleteTask(context.Background(), ""); err == nil {
		t.Fatal("expected an error")
	}
	if called {
		t.Error("no request must be sent")
	}
}

func TestDeleteTaskReturnsNotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"title":"Not Found","detail":"no such task"}`))
	}))
	defer server.Close()

	err := New(server.URL).DeleteTask(context.Background(), "d1")
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("error = %v, want *Error with 404", err)
	}
}
