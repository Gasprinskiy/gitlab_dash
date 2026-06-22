package http_client

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidBodyType = errors.New("invalid body type")
	ErrInvalidDestType = errors.New("invalid destination type")
	ErrRequestError    = errors.New("request error")
)

type RequestError struct {
	Body       []byte
	StatusCode int
}

func (re *RequestError) Error() string {
	return fmt.Sprintf("request error with status code: %d", re.StatusCode)
}
