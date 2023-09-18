package pubsubsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode string

const (
	ErrorCodeNotFound             ErrorCode = "not-found"
	ErrorCodeInvalidBody          ErrorCode = "invalid-body"
	ErrorCodeInvalidKeyReference  ErrorCode = "invalid-key-reference"
	ErrorCodeDuplicateUniqueKey   ErrorCode = "duplicate-unique-key"
	ErrorCodeCannotParseToken     ErrorCode = "cannot-parse-token"
	ErrorCodeInvalidToken         ErrorCode = "invalid-token"
	ErrorCodeMissingScopes        ErrorCode = "missing-scopes"
	ErrorCodeInternalError        ErrorCode = "internal-error"
	ErrorCodeUnauthorized         ErrorCode = "unauthorized"
	ErrorCodeInvalidAppUid        ErrorCode = "invalid-app-uid"
	ErrorCodeUnknownErrorResponse ErrorCode = "unknown-error-response"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

func newError(response *http.Response) error {
	apiError := &ErrorResponse{}

	if err := json.NewDecoder(response.Body).Decode(apiError); err != nil {
		err = fmt.Errorf("%d %s:failed to decode json error response payload: %w",
			response.StatusCode,
			http.StatusText(response.StatusCode),
			err,
		)
		return &Error{
			Code:    ErrorCodeUnknownErrorResponse,
			Message: err.Error(),
		}
	}

	return &apiError.Error
}

// Error formats the error into a string representation.
func (m *Error) Error() string {
	return fmt.Sprintf("%s: %s", m.Code, m.Message)
}
