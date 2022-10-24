package client

import "net/http"

type ResponseContext[T any] struct {
	HttpResponse *http.Response
	Data         T
}

func (r *ResponseContext[T]) StatusCode() int {
	if r.HttpResponse != nil {
		return r.HttpResponse.StatusCode
	}
	return 0
}

func (r *ResponseContext[T]) Status() string {
	if r.HttpResponse != nil {
		return r.HttpResponse.Status
	}
	return http.StatusText(0)
}
