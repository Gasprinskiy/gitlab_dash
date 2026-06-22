package http_request_proxy

import (
	"errors"
	"gitlab_api/pkg/http_client"
	"net/http"
)

var (
	ErrUnauthorizedRequest = errors.New("unauthorized request")
	ErrBadRequest          = errors.New("bad request")
	ErrForbiddenRequest    = errors.New("forbidden request")
	ErrNotFound            = errors.New("not found")
	ErrRequestTimeout      = errors.New("request timeout reached")
	ErrInternalServerError = errors.New("internal server error")
	ErrUnknownError        = errors.New("got unknown error")
)

var errCodeMap = map[int]error{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusUnauthorized:        ErrUnauthorizedRequest,
	http.StatusForbidden:           ErrForbiddenRequest,
	http.StatusRequestTimeout:      ErrRequestTimeout,
	http.StatusInternalServerError: ErrInternalServerError,
}

func HandleHttpClientRequest[T any](
	requestFunc func() (T, error),
) (T, error) {
	data, err := requestFunc()

	if err != nil {
		var (
			requestErr *http_client.RequestError
			returnErr  error
		)

		if errors.As(err, &requestErr) {
			returnErr = errCodeMap[requestErr.StatusCode]
			if returnErr == nil {
				returnErr = ErrUnknownError
			}
		} else {
			returnErr = ErrInternalServerError
		}

		return data, returnErr
	}

	return data, nil
}
