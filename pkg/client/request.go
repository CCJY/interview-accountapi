package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	DefaultContentType string = "application/json"
)

type RequestContext[T any] struct {
	HttpRequest     *http.Request
	BaseUrl         string
	OperationPath   string
	QueryParams     url.Values
	PathParams      map[string]string
	Body            interface{}
	Context         context.Context
	CustomEncoding  Encoding
	DefaultEncoding HttpEncoding
	Header          http.Header
	Method          string

	Retry         bool
	RetryInterval time.Duration
	RetryMax      int

	HttpClientDo func(*http.Request) (*http.Response, error)
}

func (r *RequestContext[T]) BuildUrl() (*url.URL, error) {
	baseUrl, err := url.Parse(r.BaseUrl)
	if err != nil {
		return nil, err
	}

	for k, v := range r.PathParams {
		r.OperationPath = strings.ReplaceAll(r.OperationPath, fmt.Sprintf("{%s}", k), v)
	}

	url, err := baseUrl.Parse(r.OperationPath)
	if err != nil {
		return nil, err
	}

	url.RawQuery = r.QueryParams.Encode()

	return url, err
}

func (r *RequestContext[T]) BuildBody() (io.Reader, error) {
	if r.Body == nil {
		return nil, nil
	}

	if r.CustomEncoding != nil {
		return r.CustomEncoding.Marshal(r.Body)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = DefaultContentType
	}
	return r.DefaultEncoding.Marshal(contentType, r.Body)
}

func (r *RequestContext[T]) NewRequest() (*RequestContext[T], error) {
	url, err := r.BuildUrl()
	if err != nil {
		return nil, err
	}

	reader, err := r.BuildBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, url.String(), reader)
	if err != nil {
		return nil, err
	}

	r.HttpRequest = req
	r.HttpRequest.Header = r.Header

	if r.Context != nil {
		r.HttpRequest = r.HttpRequest.WithContext(r.Context)
	}

	return r, err
}

func (r *RequestContext[T]) Do() (*ResponseContext[T], error) {
	req, err := r.NewRequest()
	if err != nil {
		return nil, err
	}

	rsp, err := r.HttpClientDo(req.HttpRequest)
	if err != nil {
		return nil, err
	}

	var rspData T
	if r.CustomEncoding != nil {
		err = req.CustomEncoding.UnMarshal(rsp.Body, &rspData)
		if err != nil {
			return nil, err
		}
	} else {
		err = req.DefaultEncoding.UnMarshal(rsp.Header.Get("Content-Type"), rsp.Body, &rspData)
		if err != nil {
			return nil, err
		}
	}
	defer rsp.Body.Close()

	rspContext := ResponseContext[T]{}

	rspContext.HttpResponse = rsp
	rspContext.Data = rspData

	return &rspContext, nil
}
