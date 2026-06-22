package http_client

import (
	"fmt"
	"net/http"
)

type bearerHttpTransport struct {
	transport http.RoundTripper
	token     string
}

func NewBearerHttpTransport(token string) *bearerHttpTransport {
	return &bearerHttpTransport{
		transport: http.DefaultTransport,
		token:     fmt.Sprintf("Bearer %s", token),
	}
}

func (t *bearerHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	return t.transport.RoundTrip(req)
}
