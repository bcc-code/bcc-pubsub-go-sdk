package pubsubsdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
	if gotBody["data"] != "payload" || gotBody["type"] != "my.type" {
		t.Errorf("unexpected body: %v", gotBody)
	}
	if _, ok := gotBody["message"]; ok {
		t.Error("message key must not be present")
	}
}
