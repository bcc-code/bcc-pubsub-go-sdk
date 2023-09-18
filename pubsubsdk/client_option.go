package pubsubsdk

import (
	"context"
	"net/http"
	"net/url"

	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ClientOption func(*Client)

// WithClientCredentials configures the SDK to authenticate using the client
// credentials authentication flow.
func WithClientCredentials(ctx context.Context, credEnv ClientCredentialsEnv, clientID, clientSecret string, scopes ...Scope) ClientOption {
	return func(c *Client) {
		cfg := &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     credEnv.TokenUrl,
			Scopes:       scopesToStrings(scopes),
			EndpointParams: url.Values{
				"audience": []string{credEnv.Audience},
			},
		}

		c.tokenSource = cfg.TokenSource(ctx)
	}
}

func scopesToStrings(scopes []Scope) []string {
	return lo.Map(scopes, func(s Scope, _ int) string {
		return string(s)
	})
}

// WithClientCredentials configures the SDK to authenticate using the provided token
func WithStaticToken(token string) ClientOption {
	return func(c *Client) {
		c.tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	}
}

// WithClientCredentials configures the SDK to authenticate using the provided token source
func WithTokenSource(tokenSource oauth2.TokenSource) ClientOption {
	return func(c *Client) {
		c.tokenSource = tokenSource
	}
}

// WithClient configures the SDK to use the provided client.
func WithClient(client *http.Client) ClientOption {
	return func(m *Client) {
		m.http = client
	}
}
