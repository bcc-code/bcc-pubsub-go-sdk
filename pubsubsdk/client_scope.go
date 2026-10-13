package pubsubsdk

type Scope string

const (
	ScopePubSubSubscribe Scope = "pubsub#subscribe"
	ScopePubSubPublish   Scope = "pubsub#publish"
	ScopePubSubDeliver   Scope = "pubsub#deliver"
	ScopePubSubDispatch  Scope = "pubsub#dispatch"
)
