package pubsubsdk

import (
	"context"
	"net/http"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const PublishUrl = "Publish"

func (c *Client) Publish(ctx context.Context, publishData model.PublishRequestData) error {
	return c.Request(ctx, http.MethodPost, c.URL(PublishUrl), publishData, nil)
}
