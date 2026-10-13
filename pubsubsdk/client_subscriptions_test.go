package pubsubsdk

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const subscribeResponseJSON = `{"subscriptionId":"sub1","type":"my.type","endpoint":"https://example.com/hook","clientId":"client","orgs":"a,b","scopes":"s1,s2"}`

func assertSubscribeResponse(t *testing.T, got *model.SubscribeResponse) {
	t.Helper()
	want := model.SubscribeResponse{
		SubscriptionId: "sub1",
		Type:           "my.type",
		Endpoint:       "https://example.com/hook",
		ClientId:       "client",
		Orgs:           "a,b",
		Scopes:         "s1,s2",
	}
	if got == nil || *got != want {
		t.Errorf("response = %+v, want %+v", got, want)
	}
}

func TestSubscribePostsDataAndDecodesResponse(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(subscribeResponseJSON))
	}))
	defer server.Close()

	response, err := New(server.URL).Subscribe(context.Background(), *model.NewSubscribeRequestData("my.type", "https://example.com/hook"))
	if err != nil {
		t.Fatal(err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/Subscriptions" {
		t.Errorf("path = %s, want /Subscriptions", gotPath)
	}
	if len(gotBody) != 2 || gotBody["type"] != "my.type" || gotBody["endPoint"] != "https://example.com/hook" {
		t.Errorf("unexpected body: %v", gotBody)
	}
	assertSubscribeResponse(t, response)
}

func TestSubscribeReturnsProblemDetailsError(t *testing.T) {
	server := newProblemServer(t)
	defer server.Close()

	response, err := New(server.URL).Subscribe(context.Background(), *model.NewSubscribeRequestData("my.type", "x"))
	if response != nil {
		t.Errorf("response = %+v, want nil", response)
	}
	assertBadRequestError(t, err)
}

func TestDeleteSubscription(t *testing.T) {
	tests := []struct {
		name           string
		typeName       string
		subscriptionID string
		wantURI        string
	}{
		{"without subscription id", "my.type", "", "/Subscriptions/my.type"},
		{"with subscription id", "my.type", "sub 1&x", "/Subscriptions/my.type?subscriptionId=sub+1%26x"},
		{"type with slash", "a/b", "", "/Subscriptions/a%2Fb"},
		{"type with space and question mark", "a b?c", "s", "/Subscriptions/a%20b%3Fc?subscriptionId=s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod, gotURI string
			var gotBodyLength int

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotURI = r.RequestURI
				body, _ := io.ReadAll(r.Body)
				gotBodyLength = len(body)
				_, _ = w.Write([]byte(subscribeResponseJSON))
			}))
			defer server.Close()

			response, err := New(server.URL).DeleteSubscription(context.Background(), tt.typeName, tt.subscriptionID)
			if err != nil {
				t.Fatal(err)
			}

			if gotMethod != http.MethodDelete {
				t.Errorf("method = %s, want DELETE", gotMethod)
			}
			if gotURI != tt.wantURI {
				t.Errorf("request URI = %s, want %s", gotURI, tt.wantURI)
			}
			if gotBodyLength != 0 {
				t.Errorf("body length = %d, want 0", gotBodyLength)
			}
			assertSubscribeResponse(t, response)
		})
	}
}

func TestDeleteSubscriptionEmptyTypeSendsNothing(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	response, err := New(server.URL).DeleteSubscription(context.Background(), "", "sub")
	if err == nil {
		t.Fatal("expected an error")
	}
	if response != nil {
		t.Errorf("response = %+v, want nil", response)
	}
	if called {
		t.Error("no request must be sent")
	}
}

func TestDeleteSubscriptionReturnsError(t *testing.T) {
	server := newProblemServer(t)
	defer server.Close()

	response, err := New(server.URL).DeleteSubscription(context.Background(), "my.type", "")
	if response != nil {
		t.Errorf("response = %+v, want nil", response)
	}
	assertBadRequestError(t, err)
}

func newProblemServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"title":"Bad Request","detail":"invalid input"}`))
	}))
}

func assertBadRequestError(t *testing.T, err error) {
	t.Helper()
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want *Error", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.Detail != "invalid input" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
}
