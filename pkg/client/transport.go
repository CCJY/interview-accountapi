package client

import (
	"net/http"
	"sync"
)

type Transport struct {
	Transport *http.Transport
}

type TransportOpt func(*Transport)

var DefaultTransportConfig = func() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 30
	t.MaxConnsPerHost = 30
	t.MaxIdleConnsPerHost = 30
	// t.DialContext = (&net.Dialer{
	// 	Timeout:   5 * time.Second,
	// 	KeepAlive: 5 * time.Second,
	// }).DialContext
	return t
}()

var once sync.Once

var (
	instance *Transport
)

func NewTransport(opts ...TransportOpt) *Transport {
	once.Do(func() {
		transport := &Transport{}

		for _, opt := range opts {
			opt(transport)
		}

		if transport.Transport == nil {
			transport.Transport = DefaultTransportConfig
		}

		instance = transport
	})

	return instance
}

func getTransport() *Transport {
	if instance == nil {
		return NewTransport()
	}
	return instance
}

func GetSingletonTransport() *http.Transport {
	return getTransport().Transport
}

func WithNewTransport(transport *http.Transport) TransportOpt {
	return func(t *Transport) {
		t.Transport = transport
	}
}
