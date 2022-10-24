package types

import "net/http"

type HttpResponse struct {
	HttpResponse *http.Response
}

func (r *HttpResponse) StatusCode() int {
	if r.HttpResponse != nil {
		return r.HttpResponse.StatusCode
	}
	return 0
}

func (r *HttpResponse) Status() string {
	if r.HttpResponse != nil {
		return r.HttpResponse.Status
	}
	return http.StatusText(0)
}
