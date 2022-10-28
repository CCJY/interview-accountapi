package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Retry struct {
	RetryInterval time.Duration
	RetryMax      int
}

type RetryResult struct {
	Response *http.Response
	Error    error
}

func (r *Retry) RetryRequest(client *http.Client, request *http.Request, originalBody []byte) *RetryResult {
	result := &RetryResult{
		Response: nil,
		Error:    fmt.Errorf("failed all the requests"),
	}
	ch := make(chan *RetryResult, 1)

	doFn := func(c *http.Client, req *http.Request) *RetryResult {
		got, err := c.Do(req)
		return &RetryResult{
			Response: got,
			Error:    err,
		}
	}

	// The first request depends on the timeout or context.
	// If the timeout or context is not set, wait indefinitely.
	ch <- doFn(client, request)

	// https://github.com/golang/go/issues/19653

	for retried := 0; retried < r.RetryMax; retried++ {
		select {
		case result = <-ch:
			var isError bool
			// If client.Do() has an error, response is nil
			if result.Response != nil {
				// If the http status code is 5xx, the request will be tried again.
				if http.StatusInternalServerError <= result.Response.StatusCode && result.Response.StatusCode <= 599 {
					isError = true
				}
			}

			// If err, the request will be tried again.
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
				// err "http: ContentLength=36 with Body length 0" when body has content and no error
				// buffer should be discard to reuse
				const maxBodySize = 4 << 10
				io.CopyN(io.Discard, result.Response.Body, maxBodySize)
				result.Response.Body.Close()
			}

			// request err http: ContentLength=4 with Body length 0
			// when c.Do return error and response nil
			// then reuse request
			// https://groups.google.com/g/golang-nuts/c/J-Y4LtdGNSw/m/wDSYbHWIKj0J
			// https://www.sobyte.net/post/2022-05/retry-requests/
			request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
			// unnecessary code
			// request.GetBody = func() (io.ReadCloser, error) {
			// 	return io.NopCloser(bytes.NewBuffer(originalBody)), nil
			// }

			ch <- doFn(client, request)
		case <-time.After(time.Duration(r.RetryInterval) * time.Millisecond):
		}
	}

	return result
}

func (r *Retry) Do(client *http.Client, request *http.Request, originalBody []byte) (*http.Response, error) {
	// The first request depends on the timeout or context.
	// If the timeout or context is not set, wait indefinitely.
	got, err := client.Do(request)

	result := &RetryResult{
		Response: got,
		Error:    err,
	}

	if !r.ShouldRetry(result) {
		return result.Response, result.Error
	}

	return r.Retry(client, request, originalBody)
}

func (r *Retry) ShouldRetry(result *RetryResult) bool {
	var isError bool

	// If client.Do() has an error, response is nil
	if result.Response != nil {
		// If the http status code is 5xx, the request will be tried again.
		if http.StatusInternalServerError <= result.Response.StatusCode &&
			result.Response.StatusCode <= 599 {
			isError = true
		}
	}

	// If err, the request will be tried again.
	if result.Error != nil {
		isError = true
	}

	return isError
}

func (r *Retry) Retry(client *http.Client, request *http.Request, originalBody []byte) (*http.Response, error) {
	result := &RetryResult{
		Response: nil,
		Error:    fmt.Errorf("failed all the requests"),
	}
	ch := make(chan *RetryResult, 1)

	doFn := func(c *http.Client, req *http.Request) *RetryResult {
		got, err := c.Do(req)
		return &RetryResult{
			Response: got,
			Error:    err,
		}
	}

	// The first request depends on the timeout or context.
	// If the timeout or context is not set, wait indefinitely.
	// ch <- doFn(client, request)
	ch <- func(c *http.Client, req *http.Request) *RetryResult {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.RetryInterval)*time.Millisecond)
		defer cancel()
		// Only the first retry request applies.
		newReq := req.WithContext(ctx)
		got, err := c.Do(newReq)
		return &RetryResult{
			Response: got,
			Error:    err,
		}
	}(client, request)

	// https://github.com/golang/go/issues/19653

	for retried := 1; retried <= r.RetryMax; retried++ {
		select {
		case result = <-ch:
			if !r.ShouldRetry(result) {
				return result.Response, result.Error
			}

			// Close the previous response's body. But
			// read at least some of the body so if it's
			// small the underlying TCP connection will be
			// re-used. No need to check for errors: if it
			// fails, the Transport won't reuse it anyway.
			// https://cs.opensource.google/go/go/+/master:src/net/http/client.go;l=691-695;drc=f3c39a83a3076eb560c7f687cbb35eef9b506e7d
			if result.Response != nil {
				// err "http: ContentLength=36 with Body length 0" when body has content and no error
				// buffer should be discard to reuse
				const maxBodySize = 4 << 10
				io.CopyN(io.Discard, result.Response.Body, maxBodySize)
				result.Response.Body.Close()
			}

			// request err http: ContentLength=4 with Body length 0
			// when c.Do return error and response nil
			// then reuse request
			// https://groups.google.com/g/golang-nuts/c/J-Y4LtdGNSw/m/wDSYbHWIKj0J
			// https://www.sobyte.net/post/2022-05/retry-requests/
			request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
			// unnecessary code
			// request.GetBody = func() (io.ReadCloser, error) {
			// 	return io.NopCloser(bytes.NewBuffer(originalBody)), nil
			// }

			ch <- doFn(client, request)
		case <-time.After(time.Duration(r.RetryInterval) * time.Millisecond):
		}
	}

	return result.Response, result.Error
}
