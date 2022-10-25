package client

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	BaseUrl string
	Timeout int
}

type Client struct {
	HttpClient *http.Client
	Config     ClientConfig
	// When set encoding globaly, this should set into all request context
	Encoding Encoding
}

func NewClient(t *Transport, config ClientConfig, encoding Encoding) *Client {
	return &Client{
		HttpClient: &http.Client{
			Transport: t.Transport,
			Timeout:   time.Duration(config.Timeout) * time.Second,
		},
		Config:   config,
		Encoding: encoding,
	}
}
