package http_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type HttpClient struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewHttpClient(host, token string) *HttpClient {
	base, err := url.Parse(host)
	if err != nil {
		log.Fatal("could not parse host: ", err)
	}

	return &HttpClient{
		client: &http.Client{
			Transport: NewBearerHttpTransport(token),
		},
		baseUrl: base,
	}
}

func (c *HttpClient) Get(ctx context.Context, path string) *getter {
	return &getter{
		client: c.client,
		ctx:    ctx,
		url: url.URL{
			Scheme: c.baseUrl.Scheme,
			Host:   c.baseUrl.Host,
			Path:   path,
		},
	}
}

type getter struct {
	ctx    context.Context
	query  map[string]any
	dest   any
	client *http.Client
	url    url.URL
}

func (p *getter) SetDestination(d any) *getter {
	p.dest = d
	return p
}

func (p *getter) SetQuery(q map[string]any) *getter {
	p.query = q
	return p
}

func (p *getter) Execute() (*http.Response, error) {
	if !IsValidDest(p.dest) {
		return nil, ErrInvalidDestType
	}

	if len(p.query) > 0 {
		params := url.Values{}
		for key, value := range p.query {
			params.Add(key, fmt.Sprintf("%v", value))
		}

		p.url.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(p.ctx, "GET", p.url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		body, _ := io.ReadAll(resp.Body)
		return resp, &RequestError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(p.dest); err != nil {
		return nil, err
	}

	return resp, nil
}
