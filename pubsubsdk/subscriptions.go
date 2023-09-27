package pubsubsdk

import (
	"context"
	"net/http"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const SubscriptionUrl = "Subscriptions"

func (c *Client) Subscribe(ctx context.Context, subscribeData model.SubscribeRequestData) error {
	return c.Request(ctx, http.MethodPost, c.URL(SubscriptionUrl), subscribeData, nil)
}
