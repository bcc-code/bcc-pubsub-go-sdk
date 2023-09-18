package pubsubsdk

type Scope string

const (
	ScopePubSubSubscribe Scope = "pubsub#subscribe"
	ScopePubSubPublish   Scope = "pubsub#publish"
)
