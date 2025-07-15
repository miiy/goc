package client

import (
	"net/http"
	"time"
)

type Option func(*httpClient)

func WithTimeOut(t time.Duration) Option {
	return func(c *httpClient) {
		c.Timeout = t
	}
}

func WithTransport(t *http.Transport) Option {
	return func(c *httpClient) {
		c.Transport = t
	}
}
