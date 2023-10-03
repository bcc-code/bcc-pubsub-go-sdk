package pubsubsdk

import (
	"context"
	"net/http"
	"net/url"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const SubscriptionUrl = "Subscriptions"

func (c *Client) Subscribe(ctx context.Context, subscribeData model.SubscribeRequestData) error {
	return c.Request(ctx, http.MethodPost, c.URL(SubscriptionUrl), subscribeData, nil)
}

func (c *Client) DeleteSubscription(ctx context.Context, typeName string, subscriptionId string) error {
	return c.Request(ctx, http.MethodDelete, c.URL(SubscriptionUrl+"/"+typeName+"?subscriptionId="+url.QueryEscape(subscriptionId)), nil, nil)
}
