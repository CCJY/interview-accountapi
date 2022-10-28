package client

import (
	"net/http"
	"sync"
)

// https://stuartleeks.com/posts/connection-re-use-in-golang-with-http-client/
// https://github.com/usbarmory/tamago-go/blob/c117e5d62adf00b99dc5cb9e7e0d3105d87fb09d/src/net/http/transport.go#L63-L66
// By default, Transport caches connections for future re-use.
// This may leave many open connections when accessing many hosts.
// This behavior can be managed using Transport’s CloseIdleConnections method and
// the MaxIdleConnsPerHost and DisableKeepAlives fields.
// if create new transport per one transaction,
// may be occured error for maximum socket connections
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
