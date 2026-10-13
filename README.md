# bcc-pubsub-go-sdk

Go client for the BCC PubSub API.

## Usage

Create a client with client credentials. The scopes given here are the scopes of the access token, so list every scope the client needs (`ScopePubSubPublish`, `ScopePubSubSubscribe`, `ScopePubSubDispatch`).

```go
client := pubsubsdk.New(
	pubsubsdk.UrlProd,
	pubsubsdk.WithClientCredentials(ctx, pubsubsdk.CredEnvProd, id, secret, pubsubsdk.ScopePubSubPublish),
)
```

### Publish

Needs `pubsub#publish`. `ScheduleTime` is optional and delays the publishing.

```go
data := model.NewPublishRequestData("my.type", `{"hello":"world"}`)
data.SetScheduleTime(time.Now().Add(time.Hour))
err := client.Publish(ctx, *data)
```

### Subscribe and DeleteSubscription

Need `pubsub#subscribe`. An empty subscription id refers to the default subscription of the type.

```go
subscription, err := client.Subscribe(ctx, *model.NewSubscribeRequestData("my.type", "https://example.com/hook"))

deleted, err := client.DeleteSubscription(ctx, "my.type", "")
```

### Dispatch and DeleteTask

Need `pubsub#dispatch`. `Dispatch` schedules an HTTP request to the endpoint at `ScheduleTime`.

```go
response, err := client.Dispatch(ctx, model.DispatchPayload{
	DispatchId:   "order-42",
	Type:         "my.type",
	Method:       http.MethodPost,
	Endpoint:     "https://example.com/hook",
	ScheduleTime: time.Now().Add(time.Hour),
	Data:         map[string]any{"orderId": 42},
})

err = client.DeleteTask(ctx, "order-42")
```

## Errors

Responses with status 400 or above are returned as `*pubsubsdk.Error`, with `StatusCode`, `Type`, `Title`, `Detail`, `Errors` (validation errors by field) and `Body`.

```go
var apiErr *pubsubsdk.Error
if errors.As(err, &apiErr) {
	log.Printf("status %d: %s", apiErr.StatusCode, apiErr.Detail)
}
```
