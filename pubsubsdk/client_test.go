package pubsubsdk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type recorder struct {
	srv   *httptest.Server
	count atomic.Int32
	auth  atomic.Value // string
	agent atomic.Value // string
}

func newRecorder(t *testing.T, h func(w http.ResponseWriter, r *http.Request)) *recorder {
	t.Helper()
	rec := &recorder{}
	rec.auth.Store("")
	rec.agent.Store("")
	rec.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.count.Add(1)
		rec.auth.Store(r.Header.Get("Authorization"))
		rec.agent.Store(r.Header.Get("User-Agent"))
		if h != nil {
			h(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(rec.srv.Close)
	return rec
}

func TestNewDoesNotMutateDefaultClient(t *testing.T) {
	before := http.DefaultClient.Transport
	c := New("https://example.com/pubsub", WithStaticToken("tok"))
	if http.DefaultClient.Transport != before {
		t.Fatal("http.DefaultClient.Transport was modified")
	}
	if http.DefaultClient.Timeout != 0 {
		t.Fatal("http.DefaultClient.Timeout was modified")
	}
	if c.http.Timeout != 30*time.Second {
		t.Fatalf("expected 30s timeout, got %v", c.http.Timeout)
	}
	c2 := New("https://example.com/pubsub")
	if http.DefaultClient.Transport != before {
		t.Fatal("http.DefaultClient.Transport was modified")
	}
	if c2.http.Timeout != 30*time.Second {
		t.Fatalf("expected 30s timeout, got %v", c2.http.Timeout)
	}
}

func TestNewDoesNotMutateProvidedClient(t *testing.T) {
	rt := &http.Transport{}
	custom := &http.Client{Transport: rt, Timeout: 5 * time.Second}
	c := New("https://example.com/pubsub", WithClient(custom), WithStaticToken("tok"))
	if custom.Transport != http.RoundTripper(rt) {
		t.Fatal("provided client's Transport was modified")
	}
	if custom.Timeout != 5*time.Second || c.http.Timeout != 5*time.Second {
		t.Fatalf("timeout not preserved: %v / %v", custom.Timeout, c.http.Timeout)
	}
}

func TestUserAgent(t *testing.T) {
	api := newRecorder(t, nil)
	c := New(api.srv.URL)
	if err := c.Request(context.Background(), http.MethodGet, api.srv.URL, nil, nil); err != nil {
		t.Fatal(err)
	}
	if got := api.agent.Load().(string); got != DefaultAgent {
		t.Fatalf("user agent = %q, want %q", got, DefaultAgent)
	}
}

func TestTokenSentToAPIHost(t *testing.T) {
	api := newRecorder(t, nil)
	c := New(api.srv.URL, WithStaticToken("secret"))
	if err := c.Request(context.Background(), http.MethodGet, api.srv.URL+"/x", nil, nil); err != nil {
		t.Fatal(err)
	}
	if got := api.auth.Load().(string); got != "Bearer secret" {
		t.Fatalf("Authorization = %q", got)
	}
	if got := api.agent.Load().(string); got != DefaultAgent {
		t.Fatalf("user agent = %q", got)
	}
}

func TestTokenNotSentToOtherHost(t *testing.T) {
	api := newRecorder(t, nil)
	other := newRecorder(t, nil)
	c := New(api.srv.URL, WithStaticToken("secret"))
	err := c.Request(context.Background(), http.MethodGet, other.srv.URL, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "refusing to send credentials") {
		t.Fatalf("unexpected error: %v", err)
	}
	if other.count.Load() != 0 {
		t.Fatal("other server received a request")
	}
}

func TestTokenNotFollowedOnRedirect(t *testing.T) {
	other := newRecorder(t, nil)
	api := newRecorder(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.srv.URL, http.StatusFound)
	})
	c := New(api.srv.URL, WithStaticToken("secret"))
	err := c.Request(context.Background(), http.MethodGet, api.srv.URL, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if other.count.Load() != 0 {
		t.Fatal("redirect target received a request")
	}
	if api.count.Load() != 1 {
		t.Fatalf("api server requests = %d, want 1", api.count.Load())
	}
}

func TestTokenRefusedOverInsecureHTTP(t *testing.T) {
	c := New("http://example.invalid/pubsub", WithStaticToken("secret"))
	err := c.Request(context.Background(), http.MethodGet, c.URL("x"), nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "insecure http") {
		t.Fatalf("error did not come from the guard: %v", err)
	}
}

func TestTokenRefusedWithInvalidAPIURL(t *testing.T) {
	c := New("http://[::1", WithStaticToken("secret"))
	err := c.Request(context.Background(), http.MethodGet, "https://example.com", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "invalid API URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOtherHostsAllowedWithoutToken(t *testing.T) {
	api := newRecorder(t, nil)
	other := newRecorder(t, nil)
	c := New(api.srv.URL)
	if err := c.Request(context.Background(), http.MethodGet, other.srv.URL, nil, nil); err != nil {
		t.Fatal(err)
	}
	if other.count.Load() != 1 {
		t.Fatal("other server did not receive the request")
	}
	if got := other.auth.Load().(string); got != "" {
		t.Fatalf("unexpected Authorization header %q", got)
	}
}

func TestResponseBodyLimit(t *testing.T) {
	big := newRecorder(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`"`))
		chunk := []byte(strings.Repeat("a", 1<<20))
		for i := 0; i < 11; i++ {
			_, _ = w.Write(chunk)
		}
		_, _ = w.Write([]byte(`"`))
	})
	c := New(big.srv.URL)
	var out string
	err := c.Request(context.Background(), http.MethodGet, big.srv.URL, nil, &out)
	if err == nil || !strings.Contains(err.Error(), "exceeds the limit") {
		t.Fatalf("expected limit error, got %v", err)
	}

	ok := newRecorder(t, nil)
	c = New(ok.srv.URL)
	var res struct {
		Ok bool `json:"ok"`
	}
	if err := c.Request(context.Background(), http.MethodGet, ok.srv.URL, nil, &res); err != nil {
		t.Fatal(err)
	}
	if !res.Ok {
		t.Fatal("response not decoded")
	}
}
