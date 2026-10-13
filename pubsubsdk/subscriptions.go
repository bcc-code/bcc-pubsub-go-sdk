package pubsubsdk

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const SubscriptionUrl = "Subscriptions"

// Subscribe creates or updates a subscription and returns the stored subscription.
// The credentials of the client must be granted the pubsub#subscribe scope.
func (c *Client) Subscribe(ctx context.Context, subscribeData model.SubscribeRequestData) (*model.SubscribeResponse, error) {
	var response model.SubscribeResponse
	if err := c.Request(ctx, http.MethodPost, c.URL(SubscriptionUrl), subscribeData, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteSubscription deletes the subscription for the given type and returns the deleted subscription.
// An empty subscriptionID refers to the default subscription of the type.
// The credentials of the client must be granted the pubsub#subscribe scope.
func (c *Client) DeleteSubscription(ctx context.Context, typeName string, subscriptionID string) (*model.SubscribeResponse, error) {
	if typeName == "" {
		return nil, errors.New("pubsub: subscription type must not be empty")
	}

	uri := c.URL(SubscriptionUrl, url.PathEscape(typeName))
	if subscriptionID != "" {
		uri += "?" + url.Values{"subscriptionId": {subscriptionID}}.Encode()
	}

	var response model.SubscribeResponse
	if err := c.Request(ctx, http.MethodDelete, uri, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
