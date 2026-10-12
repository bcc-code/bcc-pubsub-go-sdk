package pubsubsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxErrorBodyLength is the maximum number of bytes of the raw response body kept in Error.Body.
const maxErrorBodyLength = 1024

// Error is returned when the PubSub API responds with a status code of 400 or above.
// The API returns errors as RFC 7807 ProblemDetails, but some responses
// (for example 401/403 from auth middleware) have an empty or non-JSON body.
type Error struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"-"`
	// Type is the ProblemDetails type.
	Type string `json:"type"`
	// Title is the ProblemDetails title.
	Title string `json:"title"`
	// Detail is the ProblemDetails detail.
	Detail string `json:"detail"`
	// Body is the raw response body, truncated to 1024 bytes.
	Body string `json:"-"`
}

// newError builds an Error from an HTTP response. It always returns a non-nil *Error.
func newError(response *http.Response) *Error {
	apiError := &Error{StatusCode: response.StatusCode}

	body, _ := io.ReadAll(response.Body)

	if err := json.Unmarshal(body, apiError); err != nil {
		// Keep only the status code and the body if the payload is not ProblemDetails.
		apiError = &Error{StatusCode: response.StatusCode}
	}
	apiError.StatusCode = response.StatusCode

	if len(body) > maxErrorBodyLength {
		body = body[:maxErrorBodyLength]
	}
	apiError.Body = string(body)

	return apiError
}

// Error formats the error into a string representation.
func (e *Error) Error() string {
	message := fmt.Sprintf("pubsub: %d %s", e.StatusCode, http.StatusText(e.StatusCode))

	switch {
	case e.Detail != "":
		return message + ": " + e.Detail
	case e.Title != "":
		return message + ": " + e.Title
	case e.Body != "":
		return message + ": " + e.Body
	}

	return message
}
