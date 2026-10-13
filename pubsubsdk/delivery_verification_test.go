package pubsubsdk

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// The vectors are copied from bcc-code/bcc-pubsub (src/sdk/test-vectors/delivery-token); see
// Scripts/update-test-vectors.sh.
func TestDeliveryTokenVectors(t *testing.T) {
	files, err := filepath.Glob("testdata/delivery-token/*.json")
	if err != nil || len(files) != 14 {
		t.Fatalf("expected 14 vectors, found %d (%v)", len(files), err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			raw, _ := os.ReadFile(file)
			var vector struct {
				Issuer     string          `json:"issuer"`
				Audience   string          `json:"audience"`
				Now        int64           `json:"now"`
				JWKS       json.RawMessage `json:"jwks"`
				Token      *string         `json:"token"`
				BodyBase64 string          `json:"bodyBase64"`
				Expected   struct {
					Valid   bool   `json:"valid"`
					Subject string `json:"subject"`
					JwtID   string `json:"jwtId"`
					Reason  string `json:"reason"`
				} `json:"expected"`
			}
			if err := json.Unmarshal(raw, &vector); err != nil {
				t.Fatal(err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(vector.JWKS)
			}))
			defer server.Close()

			verifier := newTestVerifier(t, vector.Issuer, vector.Audience, server, func() time.Time { return time.Unix(vector.Now, 0) })
			body, _ := base64.StdEncoding.DecodeString(vector.BodyBase64)
			token := ""
			if vector.Token != nil {
				token = *vector.Token
			}

			claims, err := verifier.Verify(context.Background(), token, body)
			if vector.Expected.Valid {
				if err != nil {
					t.Fatalf("expected valid, got %v", err)
				}
				if claims.Subject != vector.Expected.Subject || claims.JwtID != vector.Expected.JwtID {
					t.Fatalf("claims %+v", claims)
				}
				return
			}
			if got := reasonOf(err); got != vector.Expected.Reason {
				t.Fatalf("expected %s, got %s (%v)", vector.Expected.Reason, got, err)
			}
		})
	}
}

// newTestVerifier points the verifier's key set URL at server while keeping the configured issuer.
func newTestVerifier(t *testing.T, issuer, audience string, server *httptest.Server, now func() time.Time) *DeliveryVerifier {
	t.Helper()
	verifier, err := NewDeliveryVerifier(DeliveryVerifierOptions{Issuer: issuer, Audience: audience, HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	verifier.opts.now = now
	verifier.keys.now = now
	verifier.keys.url = server.URL + "/.well-known/jwks.json"
	return verifier
}

func reasonOf(err error) string {
	var rejected *DeliveryVerificationError
	if errors.As(err, &rejected) {
		return string(rejected.Reason)
	}
	if err == nil {
		return "valid"
	}
	return err.Error()
}

const testIssuer = "https://pubsub.test"

var testBody = []byte(`{"id":"m-1"}`)

type testKey struct {
	kid string
	key *ecdsa.PrivateKey
}

func newTestKey(t *testing.T, kid string) testKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return testKey{kid: kid, key: key}
}

func (k testKey) jwk() map[string]string {
	pad := func(b []byte) string {
		out := make([]byte, 32)
		copy(out[32-len(b):], b)
		return base64.RawURLEncoding.EncodeToString(out)
	}
	return map[string]string{"kty": "EC", "crv": "P-256", "kid": k.kid, "x": pad(k.key.X.Bytes()), "y": pad(k.key.Y.Bytes())}
}

func (k testKey) sign(t *testing.T, now time.Time, override map[string]any) string {
	hash := sha256.Sum256(testBody)
	claims := map[string]any{
		"iss": testIssuer, "aud": "client-a", "sub": "s", "jti": "m-1",
		"iat": now.Unix(), "exp": now.Unix() + 300, "bdh": base64.RawURLEncoding.EncodeToString(hash[:]),
	}
	for name, value := range override {
		if value == nil {
			delete(claims, name)
		} else {
			claims[name] = value
		}
	}
	encode := func(v any) string { raw, _ := json.Marshal(v); return base64.RawURLEncoding.EncodeToString(raw) }
	signingInput := encode(map[string]string{"alg": "ES256", "typ": "JWT", "kid": k.kid}) + "." + encode(claims)
	digest := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, k.key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}

type jwksServer struct {
	mu       sync.Mutex
	keys     []testKey
	fail     bool
	requests atomic.Int32
	gate     chan struct{}
	entered  chan struct{} // when set, receives one value per request that reaches the handler
}

func (s *jwksServer) start(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		if s.entered != nil {
			s.entered <- struct{}{}
		}
		if s.gate != nil {
			<-s.gate
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.fail {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		set := map[string]any{"keys": []map[string]string{}}
		for _, k := range s.keys {
			set["keys"] = append(set["keys"].([]map[string]string), k.jwk())
		}
		_ = json.NewEncoder(w).Encode(set)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestUnknownKeyFetchesKeySetAgain(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	oldKey, newKey := newTestKey(t, "1"), newTestKey(t, "2")
	jwks := &jwksServer{keys: []testKey{oldKey}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })

	if _, err := verifier.Verify(context.Background(), oldKey.sign(t, now, nil), testBody); err != nil {
		t.Fatal(err)
	}
	now = now.Add(2 * time.Minute)
	jwks.mu.Lock()
	jwks.keys = []testKey{newKey, oldKey}
	jwks.mu.Unlock()
	if _, err := verifier.Verify(context.Background(), newKey.sign(t, now, nil), testBody); err != nil {
		t.Fatal(err)
	}
	if got := jwks.requests.Load(); got != 2 {
		t.Fatalf("expected 2 fetches, got %d", got)
	}
}

func TestUnknownKeyIsNotFetchedAgainWithinInterval(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	known := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{known}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })
	_, _ = verifier.Verify(context.Background(), known.sign(t, now, nil), testBody)

	unknown := newTestKey(t, "9").sign(t, now, nil)
	for i := 0; i < 2; i++ {
		if _, err := verifier.Verify(context.Background(), unknown, testBody); reasonOf(err) != "unknown_key" {
			t.Fatalf("expected unknown_key, got %v", err)
		}
	}
	if got := jwks.requests.Load(); got != 1 {
		t.Fatalf("expected 1 fetch, got %d", got)
	}
}

func TestConcurrentUnknownKeysFetchOnce(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	known := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{known}}
	server := jwks.start(t)
	verifier := newTestVerifier(t, testIssuer, "client-a", server, func() time.Time { return now })
	_, _ = verifier.Verify(context.Background(), known.sign(t, now, nil), testBody)
	now = now.Add(2 * time.Minute)

	const n = 20
	jwks.gate = make(chan struct{})
	jwks.entered = make(chan struct{}, n)
	unknown := newTestKey(t, "9").sign(t, now, nil)
	var wg sync.WaitGroup
	var started atomic.Int32
	results := make(chan string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			started.Add(1)
			_, err := verifier.Verify(context.Background(), unknown, testBody)
			results <- reasonOf(err)
		}()
	}
	// The first fetch is in flight and blocked; wait until every goroutine has started, so the others queue
	// behind it, then let it finish.
	<-jwks.entered
	for started.Load() < n {
		runtime.Gosched()
	}
	close(jwks.gate)
	wg.Wait()
	close(results)

	for r := range results {
		if r != "unknown_key" {
			t.Fatalf("expected unknown_key, got %s", r)
		}
	}
	if got := jwks.requests.Load(); got != 2 {
		t.Fatalf("expected 2 fetches, got %d", got)
	}
}

func TestCachedKeyIsServedDuringSlowFetch(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	known := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{known}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })
	if _, err := verifier.Verify(context.Background(), known.sign(t, now, nil), testBody); err != nil {
		t.Fatal(err)
	}
	now = now.Add(2 * time.Minute)

	jwks.gate = make(chan struct{})
	jwks.entered = make(chan struct{}, 1)
	unknownDone := make(chan struct{})
	go func() {
		defer close(unknownDone)
		_, _ = verifier.Verify(context.Background(), newTestKey(t, "9").sign(t, now, nil), testBody)
	}()
	<-jwks.entered // the fetch for the unknown kid is now blocked in the server

	cachedDone := make(chan error, 1)
	go func() {
		_, err := verifier.Verify(context.Background(), known.sign(t, now, nil), testBody)
		cachedDone <- err
	}()
	select {
	case err := <-cachedDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("a cached key waited for the fetch")
	}
	close(jwks.gate)
	<-unknownDone
}

func TestCancelledRequestDoesNotBlockRefetch(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	oldKey, newKey := newTestKey(t, "1"), newTestKey(t, "2")
	jwks := &jwksServer{keys: []testKey{oldKey}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })
	if _, err := verifier.Verify(context.Background(), oldKey.sign(t, now, nil), testBody); err != nil {
		t.Fatal(err)
	}
	now = now.Add(2 * time.Minute)
	jwks.mu.Lock()
	jwks.keys = []testKey{newKey, oldKey}
	jwks.mu.Unlock()

	jwks.gate = make(chan struct{})
	jwks.entered = make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = verifier.Verify(ctx, newKey.sign(t, now, nil), testBody)
	}()
	<-jwks.entered
	cancel() // the request that started the fetch goes away while the fetch is in flight
	close(jwks.gate)
	<-done

	if _, err := verifier.Verify(context.Background(), newKey.sign(t, now, nil), testBody); err != nil {
		t.Fatalf("expected the cancelled request not to block the key set, got %v", err)
	}
	if got := jwks.requests.Load(); got != 2 {
		t.Fatalf("expected 2 fetches, got %d", got)
	}
}

func TestCancelledContextBeforeFetchDoesNotBlockRefetch(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	key := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{key}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = verifier.Verify(ctx, key.sign(t, now, nil), testBody)

	if _, err := verifier.Verify(context.Background(), key.sign(t, now, nil), testBody); err != nil {
		t.Fatalf("expected success with a live context, got %v", err)
	}
}

func TestKeySetUnavailable(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	jwks := &jwksServer{fail: true}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })
	token := newTestKey(t, "1").sign(t, now, nil)

	for i := 0; i < 2; i++ {
		if _, err := verifier.Verify(context.Background(), token, testBody); !errors.Is(err, ErrDeliveryKeysUnavailable) {
			t.Fatalf("expected ErrDeliveryKeysUnavailable, got %v", err)
		}
	}
	if got := jwks.requests.Load(); got != 1 {
		t.Fatalf("expected 1 fetch, got %d", got)
	}
}

func TestFailedRefreshKeepsCachedKeys(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	key := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{key}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })
	_, _ = verifier.Verify(context.Background(), key.sign(t, now, nil), testBody)

	now = now.Add(2 * time.Hour)
	jwks.mu.Lock()
	jwks.fail = true
	jwks.mu.Unlock()
	if _, err := verifier.Verify(context.Background(), key.sign(t, now, nil), testBody); err != nil {
		t.Fatal(err)
	}
}

func TestMissingRequiredClaimIsMalformed(t *testing.T) {
	for _, claim := range []string{"sub", "jti", "bdh", "exp", "iat"} {
		t.Run(claim, func(t *testing.T) {
			now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
			key := newTestKey(t, "1")
			jwks := &jwksServer{keys: []testKey{key}}
			verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })

			_, err := verifier.Verify(context.Background(), key.sign(t, now, map[string]any{claim: nil}), testBody)
			if reasonOf(err) != "malformed" {
				t.Fatalf("expected malformed, got %v", err)
			}
		})
	}
}

func TestRequireDelivery(t *testing.T) {
	now := time.Now()
	key := newTestKey(t, "1")
	jwks := &jwksServer{keys: []testKey{key}}
	verifier := newTestVerifier(t, testIssuer, "client-a", jwks.start(t), func() time.Time { return now })

	var gotBody string
	var gotClaims *DeliveryClaims
	handler := RequireDelivery(verifier, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := new(strings.Builder)
		buf := make([]byte, 64)
		for {
			n, err := r.Body.Read(buf)
			raw.Write(buf[:n])
			if err != nil {
				break
			}
		}
		gotBody = raw.String()
		gotClaims, _ = DeliveryClaimsFrom(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	send := func(token, body string) int {
		req := httptest.NewRequest(http.MethodPost, "/hook", strings.NewReader(body))
		if token != "" {
			req.Header.Set(DeliveryTokenHeader, token)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := send(key.sign(t, now, nil), string(testBody)); code != http.StatusOK {
		t.Fatalf("valid push: %d", code)
	}
	if gotBody != string(testBody) || gotClaims == nil || gotClaims.JwtID != "m-1" {
		t.Fatalf("handler got body %q claims %+v", gotBody, gotClaims)
	}
	if code := send("", string(testBody)); code != http.StatusUnauthorized {
		t.Fatalf("missing token: %d", code)
	}
	if code := send(key.sign(t, now, map[string]any{"aud": "client-b"}), string(testBody)); code != http.StatusUnauthorized {
		t.Fatalf("wrong audience: %d", code)
	}
	if code := send(key.sign(t, now, nil), `{"id":"tampered"}`); code != http.StatusUnauthorized {
		t.Fatalf("changed body: %d", code)
	}
}

func TestNewDeliveryVerifierRequiresIssuerAndAudience(t *testing.T) {
	if _, err := NewDeliveryVerifier(DeliveryVerifierOptions{Audience: "a"}); err == nil {
		t.Fatal("expected error without issuer")
	}
	if _, err := NewDeliveryVerifier(DeliveryVerifierOptions{Issuer: DeliveryIssuerProd}); err == nil {
		t.Fatal("expected error without audience")
	}
}
