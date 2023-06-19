package client

import (
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	urlpkg "net/url"
	"time"
)

type Client interface {
	Get(url string, params map[string]string, headers map[string]string) (*http.Response, error)
	Post(url string, body io.Reader, headers map[string]string) (*http.Response, error)
}

type httpClient struct {
	Debug     bool
	Timeout   time.Duration
	Proxy     string
	Transport *http.Transport
	Client    *http.Client
}

var client *httpClient

func NewClient(opts ...Option) (Client, error) {
	client = &httpClient{
		Debug:   false,
		Timeout: 0,
		Proxy:   "",
		Client:  &http.Client{},
	}
	for _, o := range opts {
		o(client)
	}

	// Timeout
	// A Timeout of zero means no timeout.
	client.Client.Timeout = client.Timeout

	// cookie
	// jar, err := cookiejar.New(nil)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client.Client.Jar = jar

	return client, nil
}

func (c *httpClient) Get(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			if val != "" {
				q.Add(key, val)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	if c.Debug == true {
		u, _ := urlpkg.Parse(url)
		log.Printf("Request url: %s headers: %s cookies: %s\n", url, headers, c.Client.Jar.Cookies(u))
	}

	return c.Client.Do(req)
}

func (c *httpClient) Post(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	return c.Client.Do(req)
}
