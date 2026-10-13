package pubsubsdk

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const defaultTimeout = 30 * time.Second

type Client struct {
	http          *http.Client
	apiUrl        string
	apiAudience   string
	oauthTokenUrl string

	tokenSource oauth2.TokenSource

	httpAgent string

	// userClient is true when the caller supplied a client with WithClient.
	userClient bool
}

var DefaultAgent = fmt.Sprintf("Go-Pubsub/%s", Version)

func New(apiURL string, options ...ClientOption) *Client {
	c := &Client{
		http:      http.DefaultClient,
		apiUrl:    apiURL,
		httpAgent: DefaultAgent,
	}

	for _, option := range options {
		option(c)
	}

	if c.http == nil {
		c.http = http.DefaultClient
	}

	// Never modify the base client; it may be http.DefaultClient or owned by the caller.
	base := *c.http
	if !c.userClient {
		base.Timeout = defaultTimeout
	}
	c.http = &base

	var guard func(*http.Request) error
	if c.tokenSource != nil {
		c.http = addTokenSourceToHttpClient(c.http, c.tokenSource)
		guard = tokenGuard(apiURL)
	}

	c.http.Transport = wrappedRoundTripper{
		Base: c.http.Transport,
		Fn: func(req *http.Request) error {
			if guard != nil {
				if err := guard(req); err != nil {
					return err
				}
			}
			req.Header.Set("User-Agent", c.httpAgent)
			return nil
		},
	}

	return c
}

// tokenGuard returns a check that only lets authenticated requests through to
// the API host, over https (or http for loopback hosts).
func tokenGuard(apiURL string) func(*http.Request) error {
	parsed, parseErr := url.Parse(apiURL)
	return func(req *http.Request) error {
		if parseErr != nil {
			return fmt.Errorf("pubsubsdk: refusing to send credentials: invalid API URL: %w", parseErr)
		}
		if req.URL.Host != parsed.Host {
			return fmt.Errorf("pubsubsdk: refusing to send credentials to host %q, only %q is allowed", req.URL.Host, parsed.Host)
		}
		switch strings.ToLower(req.URL.Scheme) {
		case "https":
		case "http":
			if !isLoopbackHost(req.URL.Hostname()) {
				return errors.New("pubsubsdk: refusing to send credentials over insecure http to a non-loopback host")
			}
		default:
			return fmt.Errorf("pubsubsdk: refusing to send credentials using scheme %q", req.URL.Scheme)
		}
		return nil
	}
}

func isLoopbackHost(h string) bool {
	if strings.EqualFold(h, "localhost") {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

func addTokenSourceToHttpClient(base *http.Client, ts oauth2.TokenSource) *http.Client {
	cl := *base
	cl.Transport = &oauth2.Transport{
		Base:   base.Transport,
		Source: ts,
	}
	return &cl
}

type wrappedRoundTripper struct {
	Base http.RoundTripper
	Fn   func(*http.Request) error
}

func (w wrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	err := w.Fn(req)
	if err != nil {
		return nil, err
	}
	if w.Base != nil {
		return w.Base.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}
