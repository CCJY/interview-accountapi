package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Retry struct {
	RetryInterval time.Duration
	RetryMax      int
}

func Req(timeout int) bool {

	time.Sleep(time.Duration(timeout) * time.Millisecond)

	return true
}

func (r *Retry) TRetryRequest() (bool, bool) {
	var result bool
	ch := make(chan bool, 1)
	for i := 0; i < r.RetryMax; i++ {
		fmt.Printf("retrying:%d\n", i)
		fmt.Printf("Req timeout:%d\n", int(r.RetryInterval)-(int(r.RetryInterval)*i)/r.RetryMax)
		go func() {
			ch <- Req(int(r.RetryInterval) - (int(r.RetryInterval)*i)/r.RetryMax)
		}()
		select {
		case result = <-ch:
			return result, true
		case <-time.After(time.Duration(r.RetryInterval) * time.Millisecond):
		}

	}
	return false, false
}

type RetryResult struct {
	Response *http.Response
	Error    error
}

type RetryHttpRequestFn func(*http.Request) (*http.Response, error)

func (r *Retry) RetryRequest(client *http.Client, request *http.Request, originalBody []byte) *RetryResult {
	var result *RetryResult
	ch := make(chan *RetryResult, 1)

	doFn := func(c *http.Client, req *http.Request) *RetryResult {
		got, err := c.Do(req)
		return &RetryResult{
			Response: got,
			Error:    err,
		}
	}

	ch <- doFn(client, request)

	// https://github.com/golang/go/issues/19653

	for retried := 1; retried <= r.RetryMax; retried++ {
		fmt.Printf("Retried: %d", retried)
		select {
		case result = <-ch:
			var isError bool
			if result.Response != nil {
				if http.StatusInternalServerError <= result.Response.StatusCode && result.Response.StatusCode <= 599 {
					isError = true
				}
			}

			if result.Error != nil {
				isError = true
			}

			if !isError {
				return result
			}

			// Close the previous response's body. But
			// read at least some of the body so if it's
			// small the underlying TCP connection will be
			// re-used. No need to check for errors: if it
			// fails, the Transport won't reuse it anyway.
			// https://cs.opensource.google/go/go/+/master:src/net/http/client.go;l=691-695;drc=f3c39a83a3076eb560c7f687cbb35eef9b506e7d
			if result.Response != nil {
				// err "http: ContentLength=36 with Body length 0"
				// when body has content and no error
				// buffer should be discard to reuse
				const maxBodySize = 4 << 10
				io.CopyN(io.Discard, result.Response.Body, maxBodySize)
				result.Response.Body.Close()
			}

			// request err http: ContentLength=4 with Body length 0
			// when c.Do return error and response nil
			// then reuse request
			request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
			request.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBuffer(originalBody)), nil
			}

			ch <- doFn(client, request)
		case <-time.After(time.Duration(r.RetryInterval) * time.Millisecond):
		}
	}

	return result
}
