package pubsubsdk

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/bcc-code/bcc-pubsub-go-sdk/pubsubsdk/model"
)

const (
	DispatchUrl       = "Dispatch"
	DispatchDeleteUrl = "Delete"
)

// Dispatch schedules an HTTP request to be sent to the endpoint of the payload at its ScheduleTime.
// The credentials of the client must be granted the pubsub#dispatch scope.
func (c *Client) Dispatch(ctx context.Context, payload model.DispatchPayload) (*model.DispatchResponse, error) {
	var response model.DispatchResponse
	if err := c.Request(ctx, http.MethodPost, c.URL(DispatchUrl), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteTask deletes a task that was scheduled with Dispatch. The API responds with a
// 404 error if the task does not exist.
// The credentials of the client must be granted the pubsub#dispatch scope.
func (c *Client) DeleteTask(ctx context.Context, dispatchID string) error {
	if dispatchID == "" {
		return errors.New("pubsub: dispatch id must not be empty")
	}

	uri := c.URL(DispatchUrl, DispatchDeleteUrl) + "?" + url.Values{"dispatchId": {dispatchID}}.Encode()
	// The API answers with a plain text body, so the response is not decoded.
	return c.Request(ctx, http.MethodDelete, uri, nil, nil)
}
