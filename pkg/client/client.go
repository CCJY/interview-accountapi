package client

import (
	"net/http"
	"time"
)

type Client struct {
	HttpClient *http.Client
	Transport  *Transport
	BaseUrl    string
	// Seconds
	Timeout time.Duration
	// When set encoding globaly, this should set into all request context
	Encoding Encoding
}

type ClientOpt func(*Client)

func NewClient(opts ...ClientOpt) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	c.HttpClient = &http.Client{
		Transport: c.Transport.Transport,
		Timeout:   c.Timeout,
	}

	return c
}

func WithBaseUrl(baseUrl string) ClientOpt {
	return func(c *Client) {
		c.BaseUrl = baseUrl
	}
}

func WithTimeout(timeout int) ClientOpt {
	return func(c *Client) {
		c.Timeout = time.Duration(timeout) * time.Second
	}
}

func WithTransport(t *Transport) ClientOpt {
	return func(c *Client) {
		c.Transport = t
	}
}

func WithEncoding(encoding Encoding) ClientOpt {
	return func(c *Client) {
		c.Encoding = encoding
	}
}
