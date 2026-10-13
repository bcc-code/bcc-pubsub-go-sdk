package pubsubsdk

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DeliveryTokenHeader is the header bcc-pubsub sends the signed delivery token in.
const DeliveryTokenHeader = "X-PubSub-Token"

// Issuers of delivery tokens: the public bcc-pubsub API URL of each environment. The keys are served at
// {issuer}/.well-known/jwks.json.
const (
	DeliveryIssuerProd    = "https://pubsub.api.bcc.no"
	DeliveryIssuerSandbox = "https://pubsub.api.sandbox.bcc.no"
	DeliveryIssuerDev     = "https://pubsub.api.dev.bcc.no"
)

// DeliveryVerificationReason says why a push was rejected. The values are shared with the .NET and Node SDKs.
type DeliveryVerificationReason string

const (
	ReasonMissingToken     DeliveryVerificationReason = "missing_token"
	ReasonMalformed        DeliveryVerificationReason = "malformed"
	ReasonInvalidSignature DeliveryVerificationReason = "invalid_signature"
	ReasonUnknownKey       DeliveryVerificationReason = "unknown_key"
	ReasonWrongIssuer      DeliveryVerificationReason = "wrong_issuer"
	ReasonWrongAudience    DeliveryVerificationReason = "wrong_audience"
	ReasonExpired          DeliveryVerificationReason = "expired"
	ReasonNotYetValid      DeliveryVerificationReason = "not_yet_valid"
	ReasonBodyMismatch     DeliveryVerificationReason = "body_mismatch"
)

// DeliveryVerificationError means the push was not sent by bcc-pubsub for this client and body. Answer 401.
type DeliveryVerificationError struct {
	Reason DeliveryVerificationReason
}

func (e *DeliveryVerificationError) Error() string {
	return "delivery token rejected: " + string(e.Reason)
}

// ErrDeliveryKeysUnavailable means the issuer's key set could not be fetched. It says nothing about the token;
// answer with a 5xx status so bcc-pubsub delivers the message again.
var ErrDeliveryKeysUnavailable = errors.New("bcc-pubsub delivery keys unavailable")

// DeliveryClaims are the claims of a verified delivery token.
type DeliveryClaims struct {
	// Subject is the subscription name, or "dispatch:{dispatchId}".
	Subject string
	// JwtID is the Pub/Sub message id or dispatch id. It is the same for retries of one message.
	JwtID    string
	IssuedAt time.Time
}

// DeliveryVerifierOptions configures a DeliveryVerifier.
type DeliveryVerifierOptions struct {
	// Issuer is the public bcc-pubsub API URL, for example DeliveryIssuerProd. Required.
	Issuer string
	// Audience is your own client id. Required.
	Audience string
	// ClockSkew allowed on exp and iat. Default 60 seconds.
	ClockSkew time.Duration
	// JWKSCacheDuration is how long a fetched key set is used. Default one hour.
	JWKSCacheDuration time.Duration
	// JWKSRefreshInterval is the minimum time between fetches caused by an unknown key id. Default one minute.
	JWKSRefreshInterval time.Duration
	// HTTPClient fetches the key set. Default: a client with a 10 second timeout.
	HTTPClient *http.Client

	now func() time.Time
}

// DeliveryVerifier checks the X-PubSub-Token header bcc-pubsub sends with every push: signed by bcc-pubsub,
// issued for this client, not expired, and for exactly this body. Create one and reuse it; it caches the keys.
type DeliveryVerifier struct {
	opts   DeliveryVerifierOptions
	issuer string
	keys   *jwksCache
}

// NewDeliveryVerifier creates a verifier. It returns an error when Issuer or Audience is missing.
func NewDeliveryVerifier(opts DeliveryVerifierOptions) (*DeliveryVerifier, error) {
	if opts.Issuer == "" {
		return nil, errors.New("pubsubsdk: DeliveryVerifierOptions.Issuer is required")
	}
	if opts.Audience == "" {
		return nil, errors.New("pubsubsdk: DeliveryVerifierOptions.Audience (your client id) is required")
	}
	if opts.ClockSkew == 0 {
		opts.ClockSkew = 60 * time.Second
	}
	if opts.JWKSCacheDuration == 0 {
		opts.JWKSCacheDuration = time.Hour
	}
	if opts.JWKSRefreshInterval == 0 {
		opts.JWKSRefreshInterval = time.Minute
	}
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if opts.now == nil {
		opts.now = time.Now
	}

	issuer := strings.TrimRight(opts.Issuer, "/")
	return &DeliveryVerifier{
		opts:   opts,
		issuer: issuer,
		keys: &jwksCache{
			url:       issuer + "/.well-known/jwks.json",
			cacheFor:  opts.JWKSCacheDuration,
			refreshes: opts.JWKSRefreshInterval,
			client:    opts.HTTPClient,
			now:       opts.now,
		},
	}, nil
}

// Verify checks a token (the X-PubSub-Token header value, empty when missing) against the raw request body.
// It returns a *DeliveryVerificationError when the push must be rejected (401), or an error wrapping
// ErrDeliveryKeysUnavailable when the keys cannot be fetched (5xx).
func (v *DeliveryVerifier) Verify(ctx context.Context, token string, body []byte) (*DeliveryClaims, error) {
	if token == "" {
		return nil, reject(ReasonMissingToken)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, reject(ReasonMalformed)
	}
	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	var payload struct {
		Iss string          `json:"iss"`
		Aud json.RawMessage `json:"aud"`
		Sub string          `json:"sub"`
		Jti string          `json:"jti"`
		Bdh string          `json:"bdh"`
		Exp *float64        `json:"exp"`
		Iat *float64        `json:"iat"`
	}
	if decodeSegment(parts[0], &header) != nil || decodeSegment(parts[1], &payload) != nil {
		return nil, reject(ReasonMalformed)
	}

	// Only ES256: rejects "none" and HMAC keyed with the public key before any key is looked up.
	if header.Alg != "ES256" {
		return nil, reject(ReasonInvalidSignature)
	}
	if header.Kid == "" {
		return nil, reject(ReasonUnknownKey)
	}

	key, err := v.keys.find(ctx, header.Kid)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, reject(ReasonUnknownKey)
	}

	if !signedBy(key, parts) {
		return nil, reject(ReasonInvalidSignature)
	}

	if payload.Exp == nil || payload.Iat == nil || payload.Bdh == "" || payload.Sub == "" || payload.Jti == "" {
		return nil, reject(ReasonMalformed)
	}
	if strings.TrimRight(payload.Iss, "/") != v.issuer {
		return nil, reject(ReasonWrongIssuer)
	}
	if !hasAudience(payload.Aud, v.opts.Audience) {
		return nil, reject(ReasonWrongAudience)
	}

	now := v.opts.now()
	expires := time.Unix(int64(*payload.Exp), 0)
	issuedAt := time.Unix(int64(*payload.Iat), 0)
	if now.After(expires.Add(v.opts.ClockSkew)) {
		return nil, reject(ReasonExpired)
	}
	if issuedAt.After(now.Add(v.opts.ClockSkew)) {
		return nil, reject(ReasonNotYetValid)
	}

	hash := sha256.Sum256(body)
	if payload.Bdh != base64.RawURLEncoding.EncodeToString(hash[:]) {
		return nil, reject(ReasonBodyMismatch)
	}

	return &DeliveryClaims{Subject: payload.Sub, JwtID: payload.Jti, IssuedAt: issuedAt}, nil
}

// RequireDelivery is net/http middleware that lets a push through only with a valid X-PubSub-Token. It reads the
// body (up to maxBodyBytes; 0 means 30 MB), verifies it, and puts it back for next. A rejected push gets 401, an
// unavailable key set 503. The verified claims are in the request context (see DeliveryClaimsFrom).
func RequireDelivery(verifier *DeliveryVerifier, maxBodyBytes int64, next http.Handler) http.Handler {
	if maxBodyBytes <= 0 {
		maxBodyBytes = 30 << 20
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(DeliveryTokenHeader)
		if token == "" {
			// Decided from the headers, before the body is read.
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		claims, err := verifier.Verify(r.Context(), token, body)
		if err != nil {
			var rejected *DeliveryVerificationError
			if errors.As(err, &rejected) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), deliveryClaimsKey{}, claims)))
	})
}

type deliveryClaimsKey struct{}

// DeliveryClaimsFrom returns the claims RequireDelivery verified for this request.
func DeliveryClaimsFrom(ctx context.Context) (*DeliveryClaims, bool) {
	claims, ok := ctx.Value(deliveryClaimsKey{}).(*DeliveryClaims)
	return claims, ok
}

func reject(reason DeliveryVerificationReason) error {
	return &DeliveryVerificationError{Reason: reason}
}

func decodeSegment(segment string, target any) error {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

// signedBy checks an ES256 signature: r || s, 32 bytes each, over "header.payload".
func signedBy(key *ecdsa.PublicKey, parts []string) bool {
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || len(signature) != 64 {
		return false
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	return ecdsa.Verify(key, digest[:], r, s)
}

func hasAudience(raw json.RawMessage, audience string) bool {
	var single string
	if json.Unmarshal(raw, &single) == nil {
		return single == audience
	}
	var many []string
	if json.Unmarshal(raw, &many) == nil {
		for _, a := range many {
			if a == audience {
				return true
			}
		}
	}
	return false
}

// jwksCache holds the issuer's keys. It fetches again when the cache is older than cacheFor, or on an unknown key id
// at most once per refreshes, however many requests ask at the same time.
type jwksCache struct {
	url       string
	cacheFor  time.Duration
	refreshes time.Duration
	client    *http.Client
	now       func() time.Time

	mu        sync.Mutex
	keys      map[string]*ecdsa.PublicKey // nil until a fetch has succeeded
	fetchedAt time.Time                   // time of the last fetch attempt
	lastErr   error
}

func (c *jwksCache) find(ctx context.Context, kid string) (*ecdsa.PublicKey, error) {
	// One lock around the fetch: requests that arrive meanwhile wait for it and then use its result.
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	attempted := !c.fetchedAt.IsZero()
	if key, ok := c.keys[kid]; ok && now.Sub(c.fetchedAt) < c.cacheFor {
		return key, nil
	}
	// Fetched (or tried) a moment ago: an unknown key id does not cause another fetch.
	if attempted && now.Sub(c.fetchedAt) < c.refreshes {
		return c.lookup(kid)
	}

	keys, err := c.fetch(ctx)
	c.fetchedAt = now
	if err != nil {
		// Keep using the previous key set if there is one.
		c.lastErr = err
	} else {
		c.keys, c.lastErr = keys, nil
	}
	return c.lookup(kid)
}

func (c *jwksCache) lookup(kid string) (*ecdsa.PublicKey, error) {
	if c.keys == nil {
		return nil, fmt.Errorf("%w: could not fetch %s: %v", ErrDeliveryKeysUnavailable, c.url, c.lastErr)
	}
	return c.keys[kid], nil
}

func (c *jwksCache) fetch(ctx context.Context) (map[string]*ecdsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var set struct {
		Keys []struct {
			Kty string `json:"kty"`
			Crv string `json:"crv"`
			Kid string `json:"kid"`
			X   string `json:"x"`
			Y   string `json:"y"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&set); err != nil {
		return nil, err
	}

	keys := make(map[string]*ecdsa.PublicKey)
	for _, k := range set.Keys {
		if k.Kty != "EC" || k.Crv != "P-256" || k.Kid == "" {
			continue
		}
		if _, seen := keys[k.Kid]; seen {
			continue
		}
		if key := parseP256(k.X, k.Y); key != nil {
			keys[k.Kid] = key
		}
	}
	return keys, nil
}

// parseP256 builds a public key from JWK coordinates, or returns nil when they are not a point on P-256.
func parseP256(x, y string) *ecdsa.PublicKey {
	xb, errX := base64.RawURLEncoding.DecodeString(x)
	yb, errY := base64.RawURLEncoding.DecodeString(y)
	if errX != nil || errY != nil || len(xb) != 32 || len(yb) != 32 {
		return nil
	}
	// crypto/ecdh rejects points that are not on the curve.
	if _, err := ecdh.P256().NewPublicKey(append(append([]byte{4}, xb...), yb...)); err != nil {
		return nil
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(xb), Y: new(big.Int).SetBytes(yb)}
}
