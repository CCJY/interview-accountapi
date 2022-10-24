package client

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	HttpClient *http.Client

	// When set encoding globaly, this should set into all request context
	Encoding Encoding
}

type ClientOpt func(*Client)

var once sync.Once

var (
	instance *Client
)

func Init(opts ...ClientOpt) {
	once.Do(func() {
		client := &Client{}

		for _, opt := range opts {
			opt(client)
		}

		if client.HttpClient == nil {
			client.HttpClient = &http.Client{}
		}

		instance = client
	})

}

func GetHttpClient() *Client {
	if instance == nil {
		Init()
	}
	return instance
}

func WithNewTransport(timeSeconds int) ClientOpt {
	return func(c *Client) {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100
		c.HttpClient = &http.Client{
			Transport: t,
			Timeout:   time.Duration(timeSeconds) * time.Second,
		}
	}
}

func WithNewEncoding(encoding Encoding) ClientOpt {
	return func(c *Client) {
		c.Encoding = encoding
	}
}

type RequestInterface[T any] interface {
	Do() (*ResponseContext[T], error)
}

func NewRequestContext[T any](r *RequestContext[T]) RequestInterface[T] {
	if r != nil {
		r.HttpClientDo = GetHttpClient().HttpClient.Do

		if GetHttpClient().Encoding != nil {
			r.CustomEncoding = GetHttpClient().Encoding
		}
	}
	return r
}
