package pubsubsdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

func TestPublishPostsData(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)
	if err := client.Publish(context.Background(), *model.NewPublishRequestData("my.type", "payload")); err != nil {
		t.Fatal(err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/Publish" {
		t.Errorf("path = %s, want /Publish", gotPath)
	}
	// Only the set fields may be sent: no "message", and no empty "time" (the API rejects it).
	if len(gotBody) != 2 || gotBody["data"] != "payload" || gotBody["type"] != "my.type" {
		t.Errorf("unexpected body: %v", gotBody)
	}
}

func TestPublishSendsScheduleTime(t *testing.T) {
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
	}))
	defer server.Close()

	data := model.NewPublishRequestData("my.type", "payload")
	data.SetScheduleTime(time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC))
	if err := New(server.URL).Publish(context.Background(), *data); err != nil {
		t.Fatal(err)
	}

	if len(gotBody) != 3 || gotBody["scheduleTime"] != "2030-01-02T03:04:05Z" {
		t.Errorf("unexpected body: %v", gotBody)
	}
}
