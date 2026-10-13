package model

import "time"

// DispatchPayload describes an HTTP request that the API sends to Endpoint at ScheduleTime.
type DispatchPayload struct {
	// DispatchId identifies the task. It may only contain letters, numbers, hyphens,
	// underscores and periods, and must be 1 to 500 characters.
	DispatchId string `json:"dispatchId"`
	Type       string `json:"type"`
	// Method is the HTTP method used for the delivery, for example POST.
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	// ScheduleTime is serialized as RFC 3339.
	ScheduleTime time.Time `json:"scheduleTime"`
	// DeliverUnwrapped sends Data as the request body instead of wrapping it in an event message.
	DeliverUnwrapped bool `json:"deliverUnwrapped"`
	// Data is any JSON value. A json.RawMessage is sent as is.
	Data any `json:"data,omitempty"`
}
